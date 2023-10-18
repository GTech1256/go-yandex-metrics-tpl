package metric

import (
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
)

type memStorage struct {
	metrics *agentEntity.Metric
}

type repository struct {
	memStorage *memStorage
}

func New() *repository {
	return &repository{
		memStorage: &memStorage{metrics: nil},
	}
}
