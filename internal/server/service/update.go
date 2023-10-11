package service

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	metricvalidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"strconv"
	"strings"
)

type updateService struct {
	logger          logging2.Logger
	storage         metric2.Storage
	metricValidator metricvalidator.MetricValidator
}

func NewUpdateService(logger logging2.Logger, storage metric2.Storage, metricValidator metricvalidator.MetricValidator) updateInterface.Service {

	return &updateService{
		logger:          logger,
		storage:         storage,
		metricValidator: metricValidator,
	}
}

func (u updateService) SaveGaugeMetric(ctx context.Context, metric *entity.MetricFields) error {
	metricGaugeValue, err := u.metricValidator.GetTypeGaugeValue(metric.MetricValue)
	if err != nil {
		u.logger.Error("При получении значения метрики произошла ошибка ", err)

		return err
	}

	metricsGauge := &entity.MetricGauge{
		Type:  entity.Gauge,
		Name:  metric.MetricName,
		Value: entity.GaugeValue(*metricGaugeValue),
	}

	err = u.storage.SaveGauge(ctx, metricsGauge)
	if err != nil {
		u.logger.Error("Ошибка сохранения метрики", err)
		return err
	}

	u.logger.Info("Метрика сохранена", metricsGauge)

	return nil
}
func (u updateService) SaveCounterMetric(ctx context.Context, metric *entity.MetricFields) error {
	metricCounterValue, err := u.metricValidator.GetTypeCounterValue(metric.MetricValue)

	if err != nil {
		u.logger.Error("При получении значения метрики произошла ошибка ", err)
		return err
	}

	metricsCounter := &entity.MetricCounter{
		Type:  entity.Counter,
		Name:  metric.MetricName,
		Value: entity.CounterValue(*metricCounterValue),
	}

	err = u.storage.SaveCounter(ctx, metricsCounter)
	if err != nil {
		u.logger.Error("Ошибка сохранения метрики", err)
		return err
	}

	u.logger.Info("Метрика сохранена", metricsCounter)

	return nil
}

func (u updateService) GetMetricValue(ctx context.Context, metric *updateInterface.GetMetricValueDto) (*string, error) {
	validType := u.metricValidator.GetValidType(metric.Type)

	if validType == entity.NoType {
		return nil, metricvalidator.ErrNotCorrectType
	}

	var result *string

	switch validType {
	case entity.Counter:
		counterMetricValue, err := u.storage.GetCounterValue(metric.Name)
		if err != nil {
			u.logger.Error(err)
		}
		if counterMetricValue != nil {
			r := strconv.Itoa(int(*counterMetricValue))

			result = &r
		}
	case entity.Gauge:
		gaugeMetricValue, err := u.storage.GetGaugeValue(metric.Name)
		if err != nil {
			u.logger.Error(err)
		}
		if gaugeMetricValue != nil {
			r := fmt.Sprintf("%f", *gaugeMetricValue)

			result = &r
		}

	default:
		return nil, metricvalidator.ErrNotCorrectType
	}

	return result, nil
}

func (u updateService) GetMetrics(ctx context.Context) (string, error) {
	storageMetrics := u.storage.GetAllMetrics()
	gaugeMetrics := make([]string, len(storageMetrics.Gauge))
	counterMetrics := make([]string, len(storageMetrics.Counter))

	for name, value := range storageMetrics.Gauge {
		gaugeMetrics = append(gaugeMetrics, fmt.Sprintf("<tr>"+
			"<td>%v</td>"+
			"<td>%v</td>"+
			"</tr>", name, value),
		)
	}
	for name, value := range storageMetrics.Counter {
		counterMetrics = append(counterMetrics, fmt.Sprintf("<tr>"+
			"<td>%v</td>"+
			"<td>%v</td>"+
			"</tr>", name, value),
		)
	}

	metricsList := fmt.Sprintf("<h1>Metrics</h1>"+
		"<div style='display: flex;width: 800px;justify-content: space-between;'>"+
		"<div><h2>Gauge</h2><table><tr><th>Name</th><th>Value</th></tr>%v</table></div>"+
		"<div><h2>Counter</h2><table><tr><th>Name</th><th>Value</th></tr>%v</table></div>"+
		"</div>", strings.Join(gaugeMetrics, ""), strings.Join(counterMetrics, ""))

	html := fmt.Sprintf("%v", metricsList)

	return html, nil
}
