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
	repo := New()

	// Call the LoadMetric function
	err := repo.LoadMetric(context.Background())

	// Check for errors
	if err != nil {
		t.Fatalf("Ошибка при вызове LoadMetric: %v", err)
	}

	metric, err := repo.GetMetrics()
	if err != nil {
		t.Fatalf("Ошибка при вызове GetMetrics %v", err)
	}

	if metric == nil {
		t.Fatal("Метрики нет")
	}

	// Проверка кол-ва значений в метрике
	expectedLength := 29
	if len(*metric) != expectedLength {
		t.Fatalf("Неожиданное количество метрики. Ожидаемый: %d, Получено: %d", expectedLength, len(*metric))
	}

	for i, name := range metricFields {
		m := (*metric)[i]
		if m.MetricName != name {
			t.Errorf("Неожиданное MetricName для index %v. Ожидаемый: '%v', Получено: '%v'", i, name, m.MetricName)
		}

		expectType := string(entity.Gauge)
		if m.MetricName == "PollCount" {
			expectType = string(entity.Counter)
		}

		if m.MetricType != expectType {
			t.Errorf("Неожиданное MetricType для метрики %v. Ожидаемый: '%v', Получено: '%s'", m.MetricName, expectType, m.MetricType)
		}
	}
}
