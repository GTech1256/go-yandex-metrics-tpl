package repository

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	commonEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
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

	m := agentEntity.Metric{
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "Alloc",
			MetricValue: strconv.Itoa(int(rtm.Alloc)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "BuckHashSys",
			MetricValue: strconv.Itoa(int(rtm.BuckHashSys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "Frees",
			MetricValue: strconv.Itoa(int(rtm.Frees)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "GCCPUFraction",
			MetricValue: strconv.Itoa(int(rtm.GCCPUFraction)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "GCSys",
			MetricValue: strconv.Itoa(int(rtm.GCSys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "HeapAlloc",
			MetricValue: strconv.Itoa(int(rtm.HeapAlloc)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "HeapIdle",
			MetricValue: strconv.Itoa(int(rtm.HeapIdle)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "HeapInuse",
			MetricValue: strconv.Itoa(int(rtm.HeapInuse)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "HeapObjects",
			MetricValue: strconv.Itoa(int(rtm.HeapObjects)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "HeapReleased",
			MetricValue: strconv.Itoa(int(rtm.HeapReleased)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "HeapSys",
			MetricValue: strconv.Itoa(int(rtm.HeapSys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "LastGC",
			MetricValue: strconv.Itoa(int(rtm.LastGC)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "Lookups",
			MetricValue: strconv.Itoa(int(rtm.Lookups)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "MCacheInuse",
			MetricValue: strconv.Itoa(int(rtm.MCacheInuse)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "MCacheSys",
			MetricValue: strconv.Itoa(int(rtm.MCacheSys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "MSpanInuse",
			MetricValue: strconv.Itoa(int(rtm.MSpanInuse)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "MSpanSys",
			MetricValue: strconv.Itoa(int(rtm.MSpanSys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "Mallocs",
			MetricValue: strconv.Itoa(int(rtm.Mallocs)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "NextGC",
			MetricValue: strconv.Itoa(int(rtm.NextGC)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "NumForcedGC",
			MetricValue: strconv.Itoa(int(rtm.NumForcedGC)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "NumGC",
			MetricValue: strconv.Itoa(int(rtm.NumGC)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "OtherSys",
			MetricValue: strconv.Itoa(int(rtm.OtherSys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "PauseTotalNs",
			MetricValue: strconv.Itoa(int(rtm.PauseTotalNs)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "StackInuse",
			MetricValue: strconv.Itoa(int(rtm.StackInuse)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "StackSys",
			MetricValue: strconv.Itoa(int(rtm.StackSys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "Sys",
			MetricValue: strconv.Itoa(int(rtm.Sys)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
			MetricName:  "TotalAlloc",
			MetricValue: strconv.Itoa(int(rtm.TotalAlloc)),
		},
		{
			MetricType:  string(commonEntity.Counter),
			MetricName:  "PollCount",
			MetricValue: strconv.Itoa(int(pollCounter)),
		},
		{
			MetricType:  string(commonEntity.Gauge),
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
