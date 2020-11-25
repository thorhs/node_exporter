// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !nodiskstats

package collector

import (
	perfstat "github.com/thorhs/aix_libperfstat"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

/*
#include <sys/types.h>
#include <sys/disk.h>
*/
import "C"

type diskstatsCollector struct {
	// rxfer  typedDesc
	rbytes typedDesc
	// wxfer  typedDesc
	wbytes typedDesc
	time   typedDesc
	rtime  typedDesc
	wtime  typedDesc
	logger log.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector(logger log.Logger) (Collector, error) {
	return &diskstatsCollector{
		//rxfer:  typedDesc{readsCompletedDesc, prometheus.CounterValue},
		rbytes: typedDesc{readBytesDesc, prometheus.CounterValue},
		//wxfer:  typedDesc{writesCompletedDesc, prometheus.CounterValue},
		wbytes: typedDesc{writtenBytesDesc, prometheus.CounterValue},
		time:   typedDesc{ioTimeSecondsDesc, prometheus.CounterValue},
		rtime:  typedDesc{readTimeSecondsDesc, prometheus.CounterValue},
		wtime:  typedDesc{writeTimeSecondsDesc, prometheus.CounterValue},
		logger: logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	diskstats := perfstat.CollectDisks()

	for diskname, data := range diskstats {
		//ch <- c.rxfer.mustNewConstMetric(float64(diskstats[i].ds_rxfer), diskname)
		ch <- c.rbytes.mustNewConstMetric(data.Rblks*512, diskname)
		//ch <- c.wxfer.mustNewConstMetric(float64(diskstats[i].ds_wxfer), diskname)
		ch <- c.wbytes.mustNewConstMetric(data.Wblks*512, diskname)
		ch <- c.time.mustNewConstMetric(data.Time, diskname)
		ch <- c.rtime.mustNewConstMetric(data.Rserv, diskname)
		ch <- c.wtime.mustNewConstMetric(data.Wserv, diskname)
	}
	return nil
}
