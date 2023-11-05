package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/file"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric/converter"
	metricvalidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"strconv"
	"time"
)

type MetricValidator interface {
	GetValidType(metricType string) entity2.Type
	GetTypeGaugeValue(metricValueUnsafe string) (*float64, error)
	GetTypeCounterValue(metricValueUnsafe string) (*int64, error)
}

type Storage interface {
	SaveGauge(ctx context.Context, gauge *entity2.MetricGauge) error
	SaveCounter(ctx context.Context, counter *entity2.MetricCounter) error
	GetGaugeValue(name string) (*entity2.GaugeValue, error)
	GetCounterValue(name string) (*entity2.CounterValue, error)
	GetAllMetrics() *metric2.AllMetrics
}

type MetricLoaderService interface {
	StartMetricsToDiskInterval(ctx context.Context, interval time.Duration)
	LoadMetricsFromDisk(ctx context.Context) ([]*file.MetricJSON, error)
	SaveMetricToDisk(ctx context.Context, mj *file.MetricJSON) error
}

type metricService struct {
	logger              logging2.Logger
	storage             Storage
	metricValidator     MetricValidator
	metricLoaderService MetricLoaderService
	cfg                 *config.Config

	onMetricSave func()
}

func NewMetricService(logger logging2.Logger, storage Storage, metricValidator MetricValidator, metricLoaderService MetricLoaderService, cfg *config.Config) *metricService {
	ms := &metricService{
		logger:              logger,
		storage:             storage,
		metricValidator:     metricValidator,
		metricLoaderService: metricLoaderService,
		cfg:                 cfg,
	}

	ms.initMetricLoader()

	return ms
}

func (u metricService) SaveGaugeMetric(ctx context.Context, metric *entity2.MetricFields) error {
	metricGaugeValue, err := u.metricValidator.GetTypeGaugeValue(metric.MetricValue)
	if err != nil {
		u.logger.Error("При получении значения метрики произошла ошибка ", err)

		return err
	}

	metricsGauge := &entity2.MetricGauge{
		Type:  entity2.Gauge,
		Name:  metric.MetricName,
		Value: entity2.GaugeValue(*metricGaugeValue),
	}

	err = u.storage.SaveGauge(ctx, metricsGauge)
	if err != nil {
		u.logger.Error("Ошибка сохранения метрики", err)
		return err
	}

	u.logger.Infof("Метрика сохранена %+v", metricsGauge)

	u.saveMetricCallback(ctx, converter.MetricGaugeToMetricJSON(metricsGauge))

	return nil
}
func (u metricService) SaveCounterMetric(ctx context.Context, metric *entity2.MetricFields) error {
	metricCounterValue, err := u.metricValidator.GetTypeCounterValue(metric.MetricValue)

	if err != nil {
		u.logger.Error("При получении значения метрики произошла ошибка ", err)
		return err
	}

	metricsCounter := &entity2.MetricCounter{
		Type:  entity2.Counter,
		Name:  metric.MetricName,
		Value: entity2.CounterValue(*metricCounterValue),
	}

	err = u.storage.SaveCounter(ctx, metricsCounter)
	if err != nil {
		u.logger.Error("Ошибка сохранения метрики", err)
		return err
	}

	u.logger.Info("Метрика сохранена %+v", metricsCounter)

	u.saveMetricCallback(ctx, converter.MetricCounterToMetricJSON(metricsCounter))

	return nil
}

func (u metricService) GetMetricValue(ctx context.Context, metric *updateInterface.GetMetricValueDto) (*string, error) {
	validType := u.metricValidator.GetValidType(metric.Type)

	if validType == entity2.NoType {
		return nil, metricvalidator.ErrNotCorrectType
	}

	var result *string

	switch validType {
	case entity2.Counter:
		counterMetricValue, err := u.storage.GetCounterValue(metric.Name)
		if err != nil {
			u.logger.Error(err)
		}
		if counterMetricValue != nil {
			r := strconv.Itoa(int(*counterMetricValue))

			result = &r
		}
	case entity2.Gauge:
		gaugeMetricValue, err := u.storage.GetGaugeValue(metric.Name)
		if err != nil {
			u.logger.Error(err)
		}
		if gaugeMetricValue != nil {
			r := strconv.FormatFloat(*gaugeMetricValue, 'f', -1, 64)

			result = &r
		}

	default:
		return nil, metricvalidator.ErrNotCorrectType
	}

	return result, nil
}

func (u metricService) GetMetrics(ctx context.Context) (*metric2.AllMetrics, error) {
	storageMetrics := u.storage.GetAllMetrics()

	return storageMetrics, nil
}
