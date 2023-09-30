package updateinterface

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
)

type Service interface {
	SaveGaugeMetric(ctx context.Context, metric *entity.MetricGauge) error
	SaveCounterMetric(ctx context.Context, metric *entity.MetricCounter) error
}
