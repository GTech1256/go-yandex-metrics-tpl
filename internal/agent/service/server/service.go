package server

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"time"
)

type MetricSendCh struct {
	ID   string
	Data *agentEntity.Metric
}

type (
	Service interface {
		StartPoll(ctx context.Context, reportInterval time.Duration) error
		StartReport(ctx context.Context, reportInterval time.Duration) error
	}
)
