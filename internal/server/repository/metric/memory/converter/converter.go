package converter

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
)

func MetricJSONToMetricCounter(mj *entity2.MetricJSON) *entity.MetricCounter {
	return &entity.MetricCounter{
		Type:  entity.Counter,
		Name:  mj.ID,
		Value: *mj.Delta,
	}
}

func MetricJSONToMetricGauge(mj *entity2.MetricJSON) *entity.MetricGauge {
	return &entity.MetricGauge{
		Type:  entity.Gauge,
		Name:  mj.ID,
		Value: *mj.Value,
	}
}
