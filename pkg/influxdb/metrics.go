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

package influxdb

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ztalab/ZASentinel/pkg/influxdb/client/v2"
)

// Metrics ...
type Metrics struct {
	mu                 sync.Mutex
	conf               *CustomConfig
	batchPoints        client.BatchPoints
	point              chan *client.Point
	flushTimer         *time.Ticker
	InfluxDBHttpClient *HTTPClient
	counter            uint64
}

// MetricsData ...
type MetricsData struct {
	Measurement string                 `json:"measurement"`
	Fields      map[string]interface{} `json:"fields"`
	Tags        map[string]string      `json:"tags"`
}

// Response ...
type Response struct {
	State int      `json:"state"`
	Data  struct{} `json:"data"`
	Msg   string   `json:"msg"`
}

// NewMetrics ...
func NewMetrics(influxDBHttpClient *HTTPClient, conf *CustomConfig) (*Metrics, error) {
	bp, err := client.NewBatchPoints(influxDBHttpClient.BatchPointsConfig)
	if err != nil {
		return nil, err
	}
	metrics := &Metrics{
		conf:               conf,
		batchPoints:        bp,
		point:              make(chan *client.Point, 16),
		flushTimer:         time.NewTicker(time.Duration(conf.FlushTime) * time.Second),
		InfluxDBHttpClient: influxDBHttpClient,
	}
	go metrics.worker()
	return metrics, nil
}

func (mt *Metrics) AddPoint(metricsData *MetricsData) error {
	pt, err := client.NewPoint(metricsData.Measurement, metricsData.Tags, metricsData.Fields, time.Now())
	if err != nil {
		return err
	}
	mt.point <- pt
	return nil
}

func (mt *Metrics) worker() {
	for {
		select {
		case p, ok := <-mt.point:
			if !ok {
				mt.flush()
				return
			}
			mt.batchPoints.AddPoint(p)
			if mt.batchPoints.GetPointsNum() >= mt.conf.FlushSize {
				mt.flush()
			}
		case <-mt.flushTimer.C:
			mt.flush()
		}
	}
}

func (mt *Metrics) flush() error {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	if mt.batchPoints.GetPointsNum() == 0 {
		return nil
	}
	err := mt.Write()
	if err != nil {
		if strings.Contains(err.Error(), io.EOF.Error()) {
			err = nil
		} else {
			fmt.Println("custom-influxdb client.Write err:", err)
			return err
		}
	}
	defer mt.InfluxDBHttpClient.FluxDBHttpClose()
	mt.batchPoints.ClearPoints()
	return nil
}

func (mt *Metrics) Write() error {
	ch := make(chan error, 1)
	go func() {
		ch <- mt.InfluxDBHttpClient.FluxDBHttpWrite(mt.batchPoints)
	}()
	select {
	case err := <-ch:
		return err
	case <-time.After(800 * time.Millisecond):
		return errors.New("write timeout")
	}
}

func (mt *Metrics) count() {
	for {
		time.Sleep(time.Second)
		fmt.Println("counter：", atomic.LoadUint64(&mt.counter))
	}
}
