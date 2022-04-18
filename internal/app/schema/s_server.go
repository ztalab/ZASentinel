package schema

import (
	"github.com/ztalab/ZASentinel/pkg/errors"
	"github.com/ztalab/ZASentinel/pkg/util/json"

	jsoniter "github.com/json-iterator/go"
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
	attrByte, err := jsoniter.Marshal(attrs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = json.Unmarshal(attrByte, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
