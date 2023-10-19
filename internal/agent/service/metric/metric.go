package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

type UpdateAPI interface {
	SendUpdate(ctx context.Context, updateDto dto.Update) error
	SendUpdateJSON(ctx context.Context, updateDto dto.Update) error
}

type Repository interface {
	LoadMetric(ctx context.Context) error
	GetMetrics() (*agentEntity.Metric, error)
}

type service struct {
	server     UpdateAPI
	logger     logging.Logger
	repository Repository
}

func New(
	server UpdateAPI,
	logger logging.Logger,
	repository Repository,
) *service {

	//server.GetAPI().SendUpdate()
	return &service{
		server:     server,
		logger:     logger,
		repository: repository,
	}
}
