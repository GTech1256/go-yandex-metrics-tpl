package metric

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
)

var (
	gauge   = make(map[string]float64)
	counter = make(map[string]int64)
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

type storage struct {
	memStorage MemStorage
}

func NewStorage() metric.Storage {
	memStorage := MemStorage{
		gauge:   gauge,
		counter: counter,
	}

	return &storage{
		memStorage: memStorage,
	}
}

// SaveGauge новое значение должно замещать предыдущее.
func (s *storage) SaveGauge(ctx context.Context, gauge *entity.MetricGauge) error {
	s.memStorage.gauge[gauge.Name] = float64(gauge.Value)

	fmt.Printf("%v %+v \n", len(s.memStorage.gauge), s.memStorage)

	return nil
}

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s *storage) SaveCounter(ctx context.Context, counter *entity.MetricCounter) error {
	if _, isOk := s.memStorage.counter[counter.Name]; !isOk {
		s.memStorage.counter[counter.Name] = int64(counter.Value)
	} else {
		s.memStorage.counter[counter.Name] += int64(counter.Value)
	}

	fmt.Printf("%v %+v \n", len(s.memStorage.counter), s.memStorage)

	return nil
}

// GetGauge - возвращает значение Gauge из хранилища
func (s *storage) GetGaugeValue(name string) (entity.GaugeValue, bool) {
	value, ok := s.memStorage.gauge[name]
	return value, ok
}

// GetCounter - возвращает значение Counter из хранилища
func (s *storage) GetCounterValue(name string) (entity.CounterValue, bool) {
	value, ok := s.memStorage.counter[name]
	return value, ok
}

func (s *storage) GetAllMetrics() *metric.AllMetrics {
	return &metric.AllMetrics{
		Gauge:   s.memStorage.gauge,
		Counter: s.memStorage.counter,
	}
}
