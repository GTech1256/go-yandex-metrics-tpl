package repository

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"testing"
)

func Test_getRandomValue(t *testing.T) {
	// Проверка на разные значения
	values := make(map[uint64]bool)
	numTests := 1000

	for i := 0; i < numTests; i++ {
		randomValue := getRandomValue()
		if values[randomValue] {
			t.Errorf("Функция возвращает одинаковые значения: %d", randomValue)
		}
		values[randomValue] = true
	}
}

func Test_repository_GetMetric(t *testing.T) {
	metricFields := []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"PollCount",
		"RandomValue",
	}

	// Create an instance of the repository
	repo := repository{}

	// Call the GetMetric function
	metric, err := repo.GetMetric(context.Background())

	// Check for errors
	if err != nil {
		t.Fatalf("Error calling GetMetric: %v", err)
	}

	// Check if the returned metric is not nil
	if metric == nil {
		t.Fatal("Returned metric is nil")
	}

	// Check the length of the returned metric
	expectedLength := 29 // Adjust based on the actual number of metrics
	if len(*metric) != expectedLength {
		t.Fatalf("Unexpected number of metrics. Expected: %d, Got: %d", expectedLength, len(*metric))
	}

	for i, name := range metricFields {
		m := (*metric)[i]
		if m.MetricName != name {
			t.Errorf("Unexpected MetricName for the index %v. Expected: '%v', Got: '%v'", i, name, m.MetricName)
		}

		expectType := entity.Gauge
		if m.MetricName == "PollCount" {
			expectType = entity.Counter
		}

		if m.MetricType != expectType {
			t.Errorf("Unexpected MetricType for the metric %v. Expected: '%v', Got: '%s'", m.MetricName, expectType, m.MetricType)
		}
	}

	// Print the returned metric for reference
	t.Logf("Returned Metric: %+v", metric)
}
