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

package metrics

import (
	"context"
	"github.com/ztalab/ZASentinel/internal/config"
	"github.com/ztalab/ZASentinel/pkg/errors"
	"github.com/ztalab/ZASentinel/pkg/influxdb"
	"github.com/ztalab/ZASentinel/pkg/logger"
)

const (
	ReqSuccess = "success"
	ReqFail    = "fail"

	Prefix     = "za-sentinel"
	MetricsReq = Prefix + "req"
)

type Metrics struct {
	PodIP       string `json:"pod_ip"`
	UniqueID    string `json:"unique_id"`
	Hostname    string `json:"hostname"`
	ServiceName string `json:"service_name"`
	Delay       string `json:"delay"`
	Status      string `json:"status"`
	Operator    string `json:"operator"`
}

func AddDelayPoint(ctx context.Context, operator, status, delay, id, name string) {
	if !config.C.Influxdb.Enabled {
		return
	}
	fields := make(map[string]interface{})
	fields["delay"] = delay
	fields["status"] = status

	tags := make(map[string]string)
	tags["pod_ip"] = config.C.Common.PodIP
	tags["unique_id"] = config.C.Common.UniqueID
	tags["hostname"] = config.C.Common.Hostname
	tags["app_name"] = config.C.Common.AppName
	tags["operator"] = operator
	tags["id"] = id
	tags["name"] = name

	err := config.Is.Metrics.AddPoint(&influxdb.MetricsData{
		Measurement: MetricsReq,
		Fields:      fields,
		Tags:        tags,
	})
	if err != nil {
		logger.WithErrorStack(ctx, errors.WithStack(err)).Errorf("Failed to add sequence logs. Procedure：%v", err)
	}
}
