// Copyright 2022-present The ZTDBP Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package initer

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/ztalab/ZASentinel/pkg/certificate"
	"github.com/ztalab/ZASentinel/pkg/util"
)

const (
	TypeClient = "client"
	TypeServer = "server"
	TypeRelay  = "relay"
)

var (
	ErrCertParse = errors.New("certificate resolution error！")
	ErrCertType  = errors.New("sentinel type error！")
)

// certificate base field
type BasicCertConf struct {
	SiteID    string
	ClusterID string
	Type      string
}

func InitCert(certData []byte) (*BasicCertConf, map[string]interface{}, error) {
	p, _ := pem.Decode(certData)
	if p == nil {
		return nil, nil, ErrCertParse
	}
	cert, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, nil, ErrCertParse
	}
	basicConf := &BasicCertConf{}

	// parse attr
	mgr := certificate.New()
	attr, err := mgr.GetAttributesFromCert(cert)
	if err != nil {
		return nil, nil, ErrCertParse
	}
	if t, ok := attr.Attrs["type"]; ok {
		if t, ok := t.(string); ok {
			basicConf.Type = t
		}
	}
	if !util.InArray(basicConf.Type, []string{TypeClient, TypeRelay, TypeServer}) {
		return nil, nil, ErrCertType
	}
	return basicConf, attr.Attrs, nil
}
