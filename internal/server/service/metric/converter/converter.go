package converter

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric"
)

func MetricCounterToMetricJSON(mg *entity2.MetricCounter) *metric.MetricsJSON {
	return &metric.MetricsJSON{
		ID:    mg.Name,
		MType: string(entity.Counter),
		Delta: &mg.Value,
	}
}

func MetricGaugeToMetricJSON(mg *entity2.MetricGauge) *metric.MetricsJSON {
	return &metric.MetricsJSON{
		ID:    mg.Name,
		MType: string(entity.Gauge),
		Value: &mg.Value,
	}
}

func MetricJSONToMetricCounter(mj *metric.MetricsJSON) *entity2.MetricCounter {
	return &entity2.MetricCounter{
		Type:  entity.Counter,
		Name:  mj.ID,
		Value: *mj.Delta,
	}
}

func MetricJSONToMetricGauge(mj *metric.MetricsJSON) *entity2.MetricGauge {
	return &entity2.MetricGauge{
		Type:  entity.Gauge,
		Name:  mj.ID,
		Value: *mj.Value,
	}
}
