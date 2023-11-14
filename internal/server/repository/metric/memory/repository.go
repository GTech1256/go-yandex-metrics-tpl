package memory

import (
	"context"
	"fmt"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/sql/converter"
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

func NewStorage() *storage {
	memStorage := MemStorage{
		gauge:   gauge,
		counter: counter,
	}

	return &storage{
		memStorage: memStorage,
	}
}

// SaveGauge новое значение должно замещать предыдущее.
func (s *storage) SaveGauge(ctx context.Context, gauge *entity2.MetricGauge) error {
	s.memStorage.gauge[gauge.Name] = float64(gauge.Value)

	fmt.Printf("%v %+v \n", len(s.memStorage.gauge), s.memStorage)

	return nil
}

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s *storage) SaveCounter(ctx context.Context, counter *entity2.MetricCounter) error {
	if _, isOk := s.memStorage.counter[counter.Name]; !isOk {
		s.memStorage.counter[counter.Name] = int64(counter.Value)
	} else {
		s.memStorage.counter[counter.Name] += int64(counter.Value)
	}

	fmt.Printf("%v %+v \n", len(s.memStorage.counter), s.memStorage)

	return nil
}

// GetGauge - возвращает значение Gauge из хранилища
func (s *storage) GetGaugeValue(name string) (*entity2.GaugeValue, error) {
	value, ok := s.memStorage.gauge[name]
	if ok {
		return &value, nil
	}

	return nil, nil
}

// GetCounter - возвращает значение Counter из хранилища
func (s *storage) GetCounterValue(name string) (*entity2.CounterValue, error) {
	value, ok := s.memStorage.counter[name]
	if ok {
		return &value, nil
	}

	return nil, nil
}

func (s *storage) GetAllMetrics() *metric.AllMetrics {
	return &metric.AllMetrics{
		Gauge:   s.memStorage.gauge,
		Counter: s.memStorage.counter,
	}
}

func (s *storage) Ping(ctx context.Context) error {
	return nil
}

func (s *storage) SaveMetricBatch(ctx context.Context, metrics []*entity2.MetricJSON) error {
	for _, m := range metrics {
		if m.GetIsGauge() {
			err := s.SaveGauge(ctx, converter.MetricJSONToMetricGauge(m))
			if err != nil {
				return err
			}
		} else if m.GetIsCounter() {
			err := s.SaveCounter(ctx, converter.MetricJSONToMetricCounter(m))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
