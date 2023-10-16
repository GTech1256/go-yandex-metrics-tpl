package dto

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
)

func MetricFromService(metric *entity.MetricFields) dto.Update {
	return dto.Update{
		Type:  metric.MetricType,
		Name:  metric.MetricName,
		Value: metric.MetricValue,
	}
}
