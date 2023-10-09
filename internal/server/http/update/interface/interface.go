package updateinterface

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
)

type GetMetricValueDto struct {
	Type string
	Name string
}

type Service interface {
	SaveGaugeMetric(ctx context.Context, metric *entity.MetricFields) error
	SaveCounterMetric(ctx context.Context, metric *entity.MetricFields) error
	GetMetricValue(ctx context.Context, metric *GetMetricValueDto) (*string, error)
	GetMetrics(ctx context.Context) (string, error)
}
