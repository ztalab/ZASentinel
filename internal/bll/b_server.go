// Copyright 2022-present The Ztalab Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bll

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/xtaci/smux"
	"github.com/ztalab/ZASentinel/internal/config"
	"github.com/ztalab/ZASentinel/internal/contextx"
	"github.com/ztalab/ZASentinel/internal/event"
	"github.com/ztalab/ZASentinel/internal/metrics"
	"github.com/ztalab/ZASentinel/internal/schema"
	"github.com/ztalab/ZASentinel/pkg/certificate"
	"github.com/ztalab/ZASentinel/pkg/errors"
	"github.com/ztalab/ZASentinel/pkg/logger"
	"github.com/ztalab/ZASentinel/pkg/pconst"
	"github.com/ztalab/ZASentinel/pkg/recover"
	"github.com/ztalab/ZASentinel/pkg/util/json"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type Server struct{}

// ReadInitiaWSRequest Read the WS request
func (a *Server) ReadInitiaWSRequest(ctx context.Context, connReader *bufio.Reader, conf *schema.ServerConfig) (*schema.ClientConfig, *http.Request, context.Context, error) {
	expectedH1Req := "GET /secretLink"
	firstBytes, err := connReader.Peek(len(expectedH1Req))
	if err != nil {
		return nil, nil, ctx, errors.WithStack(err)
	}
	if string(firstBytes) == expectedH1Req {
		req, err := http.ReadRequest(connReader)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		traceID := req.Header.Get("X-TraceID")
		if traceID != "" {
			ctx = contextx.NewTraceID(ctx, traceID)
			ctx = logger.NewTraceIDContext(ctx, traceID)
		}
		if strings.ToLower(req.Header.Get("Connection")) != "upgrade" && strings.ToLower(req.Header.Get("Connection")) != "keep-alive, upgrade" {
			return nil, nil, ctx, errors.WithStack(fmt.Errorf("Connection header expected: upgrade, got: %s\n",
				strings.ToLower(req.Header.Get("Connection"))))
		}
		if strings.ToLower(req.Header.Get("Upgrade")) != "websocket" {
			return nil, nil, ctx, errors.WithStack(fmt.Errorf("Upgrade header expected: websocket, got: %s\n",
				strings.ToLower(req.Header.Get("Upgrade"))))
		}
		// Get link information
		chainsJSON := strings.ToLower(req.Header.Get("X-Chains"))
		if chainsJSON == "" {
			return nil, nil, ctx, errors.NewWithStack("X-Chains argument is missing")
		}
		var chains schema.ClientConfig
		err = json.Unmarshal([]byte(chainsJSON), &chains)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		// Verify the resources
		if ok := conf.Resources.VerifyResources(chains.Target); !ok {
			err := errors.New("The server verifies that the requested resource does not exist")
			event.NewServerEvent(&chains, conf, event.TagResourceNotFound, err.Error()).Error(ctx)
			return nil, nil, ctx, errors.WithStack(err)
		}
		// get client certificate
		clientCa := req.Header.Get("X-ClientCert")
		if clientCa == "" {
			return nil, nil, ctx, errors.NewWithStack("X-ClientCert argument is missing")
		}
		clientCaCert, err := base64.StdEncoding.DecodeString(clientCa)
		if err != nil {
			return nil, nil, ctx, errors.WithStack(err)
		}
		// Verify the client certificate
		err = certificate.NewVerify(string(clientCaCert), config.C.Certificate.CaPem, "").Verify()
		if err != nil {
			event.NewServerEvent(&chains, conf, event.TagClientTLSFail, err.Error()).Error(ctx)
			return nil, nil, ctx, errors.WithStack(err)
		}
		return &chains, req, ctx, nil
	}
	req, err := http.ReadRequest(connReader)
	if err != nil {
		return nil, nil, ctx, errors.WithStack(err)
	}
	reqBytes, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, nil, ctx, errors.WithStack(err)
	}
	err = errors.New("Server Illegal request:\n" + string(reqBytes))
	return nil, nil, ctx, errors.WithStack(err)
}

// Responding to WS requests
func (a *Server) GenerateInitialWSResponse(ctx context.Context, clientConn net.Conn, req *http.Request) ([]byte, error) {
	resp := http.Response{
		Status:           "101 Switching Protocols",
		StatusCode:       101,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           http.Header{},
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
	resp.Header.Set("Upgrade", req.Header.Get("Upgrade"))
	resp.Header.Set("Connection", req.Header.Get("Connection"))
	resp.Header.Set("X-ServerCert", base64.StdEncoding.EncodeToString([]byte(config.C.Certificate.CertPem)))

	res, err := httputil.DumpResponse(&resp, true)
	_, err = clientConn.Write(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, err
}

func (a *Server) handleConn(ctx context.Context, conf *schema.ServerConfig, clientConn net.Conn) error {
	begin := time.Now()
	connReader := bufio.NewReader(clientConn)
	chains, req, ctx, err := a.ReadInitiaWSRequest(ctx, connReader, conf)
	if err != nil {
		metrics.AddDelayPoint(ctx, pconst.OperatorServer, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
		logger.WithErrorStack(ctx, err).Error("Error obtaining WS request information：", err)
		return err
	}
	_, err = a.GenerateInitialWSResponse(ctx, clientConn, req)
	if err != nil {
		metrics.AddDelayPoint(ctx, pconst.OperatorServer, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
		logger.WithErrorStack(ctx, err).Error("Response WS message error：", err)
		return err
	}
	verifyFlag := "serverCaReady"
	verifyBytes, err := connReader.Peek(len(verifyFlag))
	if err != nil {
		metrics.AddDelayPoint(ctx, pconst.OperatorServer, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
		logger.WithErrorStack(ctx, err).Error("Error obtaining the certificate verification result.：", err)
		return err
	}
	if string(verifyBytes) == verifyFlag {
		targetAddr := chains.Target.Host + ":" + strconv.Itoa(chains.Target.Port)
		serverConn, err := net.Dial("tcp", targetAddr)
		if err != nil {
			metrics.AddDelayPoint(ctx, pconst.OperatorServer, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
			event.NewServerEvent(chains, conf, event.TagConnectFail, err.Error()).Error(ctx)
			logger.WithErrorStack(ctx, errors.WithStack(err)).Errorf("Failed to request resource from server\n:Addr:%s Error:%v", targetAddr, err)
			return err
		}
		metrics.AddDelayPoint(ctx, pconst.OperatorServer, metrics.ReqSuccess, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
		event.NewServerEvent(chains, conf, event.TagConnectSuccess, "").Info(ctx)
		// 多路复用
		session, err := smux.Server(clientConn, nil)
		if err != nil {
			return err
		}
		stream, err := session.AcceptStream()
		if err != nil {
			return err
		}
		defer func() {
			closeErr := clientConn.Close()
			if closeErr != nil {
				logger.WithContext(ctx).Errorf("Closed Connection with error: %v\n", closeErr)
			} else {
				logger.WithContext(ctx).Infof("Closed Connection: %v\n", clientConn.RemoteAddr().String())
			}
			serverConn.Close()
			session.Close()
			stream.Close()
		}()
		TransparentProxy(stream, serverConn)
		return nil
	}
	err = errors.New("Server certificate verification failed")
	logger.WithErrorStack(ctx, errors.WithStack(err)).Error(err)
	metrics.AddDelayPoint(ctx, pconst.OperatorServer, metrics.ReqFail, time.Now().Sub(begin).String(), conf.UUID, conf.Name)
	return err
}

func NewServer() *Server {
	return &Server{}
}

func (a *Server) Listen(ctx context.Context, attrs map[string]interface{}) {
	go func() {
		conf, err := schema.ParseServerConfig(attrs)
		if err != nil {
			panic(err)
		}
		cert, err := tls.X509KeyPair([]byte(config.C.Certificate.CertPem), []byte(config.C.Certificate.KeyPem))
		if err != nil {
			panic(err)
		}
		l, err := tls.Listen("tcp", "0.0.0.0:"+strconv.Itoa(conf.Port), &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			panic(err)
		}
		logger.WithContext(ctx).Printf("Started ZERO ACCESS Server at %v\n", l.Addr().String())
		for {
			conn, err := l.Accept()
			if err != nil {
				logger.WithContext(ctx).Error("Failed to accept connection:", err)
				continue
			}
			recover.Recovery(ctx, func() {
				a.handleConn(ctx, conf, conn)
			})
		}
	}()
}
