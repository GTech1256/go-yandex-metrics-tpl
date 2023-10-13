package server

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
)

type Repository interface {
	LoadMetric(ctx context.Context) error
	GetMetrics() (*agentEntity.Metric, error)
}
