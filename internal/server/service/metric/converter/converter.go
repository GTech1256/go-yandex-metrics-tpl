package converter

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/file"
)

func MetricCounterToMetricJSON(mg *entity.MetricCounter) *file.MetricJSON {
	return &file.MetricJSON{
		ID:    mg.Name,
		MType: string(entity.Counter),
		Delta: &mg.Value,
	}
}

func MetricGaugeToMetricJSON(mg *entity.MetricGauge) *file.MetricJSON {
	return &file.MetricJSON{
		ID:    mg.Name,
		MType: string(entity.Gauge),
		Value: &mg.Value,
	}
}

func MetricJSONToMetricCounter(mj *file.MetricJSON) *entity.MetricCounter {
	return &entity.MetricCounter{
		Type:  entity.Counter,
		Name:  mj.ID,
		Value: *mj.Delta,
	}
}

func MetricJSONToMetricGauge(mj *file.MetricJSON) *entity.MetricGauge {
	return &entity.MetricGauge{
		Type:  entity.Gauge,
		Name:  mj.ID,
		Value: *mj.Value,
	}
}
