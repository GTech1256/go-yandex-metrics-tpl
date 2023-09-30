package updateinterface

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
)

type Service interface {
	GetMetric(ctx context.Context, url string) (*entity.Metric, error)
	SaveGaugeMetric(ctx context.Context, metric *entity.MetricGauge) error
	SaveCounterMetric(ctx context.Context, metric *entity.MetricCounter) error
}
