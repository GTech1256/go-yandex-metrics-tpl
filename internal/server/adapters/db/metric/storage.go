package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
)

var (
	gauge   = make(map[string]float64)
	counter = make(map[string]int64)
)

type MemStorage = struct {
	gauge   map[string]float64
	counter map[string]int64
}

var MemStorageContainer MemStorage = MemStorage{
	gauge:   gauge,
	counter: counter,
}

type storage struct {
}

func NewStorage() metric.Storage {
	return &storage{}
}

// SaveGauge новое значение должно замещать предыдущее.
func (s storage) SaveGauge(ctx context.Context, gauge *entity.MetricGauge) error {
	MemStorageContainer.gauge[gauge.Name] = float64(gauge.Value)

	return nil
}

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s storage) SaveCounter(ctx context.Context, counter *entity.MetricCounter) error {
	if _, isOk := gauge[counter.Name]; !isOk {
		MemStorageContainer.counter[counter.Name] = int64(counter.Value)
	} else {
		MemStorageContainer.counter[counter.Name] += int64(counter.Value)
	}

	return nil
}
