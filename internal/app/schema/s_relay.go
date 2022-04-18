package schema

import (
	"github.com/ztalab/ZASentinel/pkg/errors"
	"github.com/ztalab/ZASentinel/pkg/util/json"

	jsoniter "github.com/json-iterator/go"
)

type RelayConfig struct {
	Type string `json:"type"`
	Port int    `json:"port"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func ParseRelayConfig(attrs map[string]interface{}) (*RelayConfig, error) {
	var result RelayConfig
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
