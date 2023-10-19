package converter

import (
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/models"
)

func MetricsToMetricValueDTO(metrics models.Metrics) updateInterface.GetMetricValueDto {
	return updateInterface.GetMetricValueDto{
		Type: metrics.MType,
		Name: metrics.ID,
	}
}
