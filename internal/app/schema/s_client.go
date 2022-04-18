package schema

import (
	"sort"
	"github.com/ztalab/ZASentinel/pkg/errors"
	"github.com/ztalab/ZASentinel/pkg/util/json"

	jsoniter "github.com/json-iterator/go"
)

type ClientConfig struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Port      int       `json:"port"`
	Relays    Relays    `json:"relay"`
	Server    Server    `json:"server"`
	Target    Target    `json:"target"`
	Resources Resources `json:"resources"`
}

type Relay struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	OutPort int    `json:"out_port"`
	Sort    int    `json:"sort"`
}

type Server struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	OutPort int    `json:"out_port"`
}

type Target struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CResource struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type (
	Relays     []*Relay
	CResources []*CResource
)

func ParseClientConfig(attrs map[string]interface{}) (*ClientConfig, error) {
	var result ClientConfig
	attrByte, err := jsoniter.Marshal(attrs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = json.Unmarshal(attrByte, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(result.Relays) > 0 {
		// Reply sort
		result.RelaysAscBySort()
	}
	if result.Server.Host == "" {
		err := errors.New("server Addr argument is missing")
		return nil, errors.WithStack(err)
	}
	return &result, nil
}

func (a *ClientConfig) ToJSONString() string {
	return json.MarshalToString(a)
}

// RelaysAscBySort
func (a *ClientConfig) RelaysAscBySort() {
	sort.Slice(a.Relays, func(i, j int) bool { // asc
		return a.Relays[i].Sort < a.Relays[j].Sort
	})
}
