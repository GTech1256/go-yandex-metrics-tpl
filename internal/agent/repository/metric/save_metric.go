package metric

import agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"

func (r *repository) saveMetrics(metrics *agentEntity.Metric) error {
	r.memStorage.metrics = metrics

	return nil
}
