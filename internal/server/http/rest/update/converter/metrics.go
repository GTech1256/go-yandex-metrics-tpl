package converter

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"strconv"
)

func MetricsGaugeToMetricFields(metrics entity.MetricsJSON) entity.MetricFields {
	return entity.MetricFields{
		MetricType:  metrics.MType,
		MetricName:  metrics.ID,
		MetricValue: strconv.FormatFloat(*metrics.Value, 'f', -1, 64),
	}
}
func MetricsCounterToMetricFields(metrics entity.MetricsJSON) entity.MetricFields {
	return entity.MetricFields{
		MetricType:  metrics.MType,
		MetricName:  metrics.ID,
		MetricValue: strconv.Itoa(int(*metrics.Delta)),
	}
}

func MetricsToMetricValueDTO(metrics entity.MetricsJSON) updateInterface.GetMetricValueDto {
	return updateInterface.GetMetricValueDto{
		Type: metrics.MType,
		Name: metrics.ID,
	}
}
