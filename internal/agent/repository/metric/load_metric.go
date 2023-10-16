package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var (
	pollCounter = 0
	rtm         runtime.MemStats
)

func (r *repository) LoadMetric(ctx context.Context) error {
	RandomValue := getRandomValue()
	pollCounter++

	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	m := entity.Metric{
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "Alloc",
			MetricValue: strconv.Itoa(int(rtm.Alloc)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "BuckHashSys",
			MetricValue: strconv.Itoa(int(rtm.BuckHashSys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "Frees",
			MetricValue: strconv.Itoa(int(rtm.Frees)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "GCCPUFraction",
			MetricValue: strconv.Itoa(int(rtm.GCCPUFraction)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "GCSys",
			MetricValue: strconv.Itoa(int(rtm.GCSys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "HeapAlloc",
			MetricValue: strconv.Itoa(int(rtm.HeapAlloc)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "HeapIdle",
			MetricValue: strconv.Itoa(int(rtm.HeapIdle)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "HeapInuse",
			MetricValue: strconv.Itoa(int(rtm.HeapInuse)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "HeapObjects",
			MetricValue: strconv.Itoa(int(rtm.HeapObjects)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "HeapReleased",
			MetricValue: strconv.Itoa(int(rtm.HeapReleased)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "HeapSys",
			MetricValue: strconv.Itoa(int(rtm.HeapSys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "LastGC",
			MetricValue: strconv.Itoa(int(rtm.LastGC)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "Lookups",
			MetricValue: strconv.Itoa(int(rtm.Lookups)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "MCacheInuse",
			MetricValue: strconv.Itoa(int(rtm.MCacheInuse)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "MCacheSys",
			MetricValue: strconv.Itoa(int(rtm.MCacheSys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "MSpanInuse",
			MetricValue: strconv.Itoa(int(rtm.MSpanInuse)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "MSpanSys",
			MetricValue: strconv.Itoa(int(rtm.MSpanSys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "Mallocs",
			MetricValue: strconv.Itoa(int(rtm.Mallocs)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "NextGC",
			MetricValue: strconv.Itoa(int(rtm.NextGC)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "NumForcedGC",
			MetricValue: strconv.Itoa(int(rtm.NumForcedGC)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "NumGC",
			MetricValue: strconv.Itoa(int(rtm.NumGC)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "OtherSys",
			MetricValue: strconv.Itoa(int(rtm.OtherSys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "PauseTotalNs",
			MetricValue: strconv.Itoa(int(rtm.PauseTotalNs)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "StackInuse",
			MetricValue: strconv.Itoa(int(rtm.StackInuse)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "StackSys",
			MetricValue: strconv.Itoa(int(rtm.StackSys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "Sys",
			MetricValue: strconv.Itoa(int(rtm.Sys)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "TotalAlloc",
			MetricValue: strconv.Itoa(int(rtm.TotalAlloc)),
		},
		{
			MetricType:  string(entity.Counter),
			MetricName:  "PollCount",
			MetricValue: strconv.Itoa(int(pollCounter)),
		},
		{
			MetricType:  string(entity.Gauge),
			MetricName:  "RandomValue",
			MetricValue: strconv.Itoa(int(RandomValue)),
		},
	}

	return r.saveMetrics(&m)
}

func getRandomValue() uint64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	RandomValue := rand.Uint64()

	return RandomValue
}
