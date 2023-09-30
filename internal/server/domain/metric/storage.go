package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
)

type Storage interface {
	SaveGauge(ctx context.Context, gauge *entity.MetricGauge) error
	SaveCounter(ctx context.Context, counter *entity.MetricCounter) error
}
