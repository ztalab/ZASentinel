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

package event

import (
	"context"
	"github.com/ztalab/ZASentinel/internal/schema"
	"github.com/ztalab/ZASentinel/pkg/logger"
	"github.com/ztalab/ZASentinel/pkg/pconst"
	"github.com/ztalab/ZASentinel/pkg/util/json"

	"github.com/sirupsen/logrus"
)

const (
	TagConnectSuccess   = "Connect success"
	TagConnectFail      = "Connect fail"
	TagClientTLSFail    = "Client tls invalid"
	TagServerTLSFail    = "Server tls invalid"
	TagResourceNotFound = "Resource not found"
)

type Event struct {
	Operator   string               `json:"operator"`
	ClientInfo *schema.ClientConfig `json:"client_info"`
	ServerInfo *schema.ServerConfig `json:"server_info"`
	RelayInfo  *schema.RelayConfig  `json:"relay_info"`
	Tag        string               `json:"tag"`
	MsgInfo    string               `json:"msg_info"`
}

func NewClientEvent(clientInfo *schema.ClientConfig, tag, msgInfo string) *Event {
	return &Event{
		Operator:   pconst.OperatorClient,
		ClientInfo: clientInfo,
		Tag:        tag,
		MsgInfo:    msgInfo,
	}
}

func NewRelayEvent(clientInfo *schema.ClientConfig, relayInfo *schema.RelayConfig, tag, msgInfo string) *Event {
	return &Event{
		Operator:   pconst.OperatorRelay,
		ClientInfo: clientInfo,
		RelayInfo:  relayInfo,
		Tag:        tag,
		MsgInfo:    msgInfo,
	}
}

func NewServerEvent(clientInfo *schema.ClientConfig, serverInfo *schema.ServerConfig, tag, msgInfo string) *Event {
	return &Event{
		Operator:   pconst.OperatorServer,
		ClientInfo: clientInfo,
		ServerInfo: serverInfo,
		Tag:        tag,
		MsgInfo:    msgInfo,
	}
}

func (a *Event) Info(ctx context.Context) {
	a.toLog(ctx, logrus.InfoLevel)
}

func (a *Event) Error(ctx context.Context) {
	a.toLog(ctx, logrus.ErrorLevel)
}

func (a *Event) Warn(ctx context.Context) {
	a.toLog(ctx, logrus.WarnLevel)
}

func (a *Event) toLog(ctx context.Context, level logrus.Level) {
	msg := json.MarshalToString(a)
	fields := make(map[string]interface{})
	fields["customLog1"] = "Events"
	fields["customLog1"] = a.Operator
	switch level {
	case logrus.ErrorLevel:
		logger.WithContext(ctx).WithFields(fields).Errorf("Events：%s", msg)
	case logrus.WarnLevel:
		logger.WithContext(ctx).WithFields(fields).Warnf("Events：%s", msg)
	case logrus.InfoLevel:
		logger.WithContext(ctx).WithFields(fields).Infof("Events：%s", msg)
	}
}
