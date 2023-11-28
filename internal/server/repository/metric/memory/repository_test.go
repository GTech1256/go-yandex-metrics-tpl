package memory

import (
	"context"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"testing"
)

func TestSaveGauge(t *testing.T) {
	storage := NewStorage()

	gauge := &entity2.MetricGauge{
		Type:  entity2.Gauge,
		Name:  "test_gauge",
		Value: 42.0,
	}

	err := storage.SaveGauge(context.Background(), gauge)
	if err != nil {
		t.Fatalf("SaveGauge failed: %v", err)
	}

	// Проверяем, что значение было сохранено
	if storedValue, err2 := storage.GetGaugeValue(gauge.Name); err2 != nil || *storedValue != float64(gauge.Value) {
		t.Errorf("SaveGauge did not save the value correctly")
	}
}

func TestSaveCounter(t *testing.T) {
	storage := NewStorage()

	counter := &entity2.MetricCounter{
		Type:  entity2.Gauge,
		Name:  "test_counter",
		Value: 10,
	}

	// Тест случая, когда значение не было известно
	err := storage.SaveCounter(context.Background(), counter)
	if err != nil {
		t.Fatalf("SaveCounter failed: %v", err)
	}

	// Проверяем, что значение было сохранено и равно переданному
	if storedValue, err2 := storage.GetCounterValue(counter.Name); err2 != nil || *storedValue != int64(counter.Value) {
		t.Errorf("SaveCounter did not save the value correctly")
	}

	// Тест случая, когда значение уже было известно
	previousValue := int64(counter.Value) // начинаем с текущего значения счетчика
	err = storage.SaveCounter(context.Background(), counter)
	if err != nil {
		t.Fatalf("SaveCounter failed: %v", err)
	}

	err = storage.SaveCounter(context.Background(), counter)
	if err != nil {
		t.Fatalf("SaveCounter failed: %v", err)
	}

	// Проверяем, что значение увеличено на переданное
	expectedValue := previousValue + 2*int64(counter.Value)
	if storedValue, err2 := storage.GetCounterValue(counter.Name); err2 != nil || *storedValue != expectedValue {
		t.Errorf("SaveCounter did not update the value correctly Expect: %v, Got: %v", expectedValue, storedValue)
	}
}
