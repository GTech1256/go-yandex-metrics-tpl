package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
)

type AllMetrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type Storage interface {
	SaveGauge(ctx context.Context, gauge *entity.MetricGauge) error
	SaveCounter(ctx context.Context, counter *entity.MetricCounter) error
	GetGaugeValue(name string) (entity.GaugeValue, bool)
	GetCounterValue(name string) (entity.CounterValue, bool)
	GetAllMetrics() *AllMetrics
}
