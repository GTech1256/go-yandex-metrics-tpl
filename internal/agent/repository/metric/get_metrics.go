package metric

import agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"

func (r *repository) GetMetrics() (*agentEntity.Metrics, error) {
	return r.memStorage.metrics, nil
}
