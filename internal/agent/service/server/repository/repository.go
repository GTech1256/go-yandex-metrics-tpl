package repository

import (
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
)

type MemStorage struct {
	metrics *agentEntity.Metric
}

type repository struct {
	memStorage *MemStorage
}

func New() server.Repository {
	return &repository{
		memStorage: &MemStorage{metrics: nil},
	}
}

func (r *repository) saveMetrics(metrics *agentEntity.Metric) error {
	r.memStorage.metrics = metrics

	return nil
}

func (r *repository) GetMetrics() (*agentEntity.Metric, error) {
	return r.memStorage.metrics, nil
}
