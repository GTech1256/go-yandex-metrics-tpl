package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/config"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"sync"
)

type UpdateAPI interface {
	SendUpdate(ctx context.Context, updateDto dto.Update) error
	SendUpdateJSON(ctx context.Context, updateDto dto.Update) error
	SendUpdates(ctx context.Context, updateDto []*dto.Update) error
}

type Repository interface {
	LoadMetric(ctx context.Context) error
	GetMetrics() (*entity.Metrics, error)
}

type service struct {
	server     UpdateAPI
	logger     logging.Logger
	repository Repository
	cfg        *config.Config

	// В идеале в use-case должен быть, а не полностью в сервисе
	// Канал, который содержит метрики для отправки для RateLimit
	metricsForSendCh chan *entity.MetricFields

	metricWorkerOnce sync.Once
}

func New(
	server UpdateAPI,
	logger logging.Logger,
	repository Repository,
	cfg *config.Config,

) *service {

	return &service{
		server:     server,
		logger:     logger,
		repository: repository,
		cfg:        cfg,

		// Канал, который содержит метрики для отправки через воркер
		metricsForSendCh: make(chan *entity.MetricFields, 1024),
	}
}
