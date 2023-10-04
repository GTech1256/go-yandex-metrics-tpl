package server

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"time"
)

type (
	Service interface {
		StartPoll(ctx context.Context, metricSendCh chan<- *agentEntity.Metric, reportInterval time.Duration) error
		StartReport(ctx context.Context, metricSendCh <-chan *agentEntity.Metric, pollInterval time.Duration) error
		SendMetric(ctx context.Context, metric *entity.MetricFields) error
	}
)
