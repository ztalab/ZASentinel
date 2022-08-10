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

package schema

import (
	"github.com/ztalab/ZASentinel/pkg/errors"
	"github.com/ztalab/ZASentinel/pkg/util/json"
)

type ServerConfig struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Port      int       `json:"port"`
	Resources Resources `json:"resource"`
}

func ParseServerConfig(attrs map[string]interface{}) (*ServerConfig, error) {
	var result ServerConfig
	attrByte, err := json.Marshal(attrs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = json.Unmarshal(attrByte, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
