package metric

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/file"
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
	LoadMetricsFromDisk(ctx context.Context) ([]*file.MetricsJSON, error)
	SaveMetricToDisk(ctx context.Context, mj *file.MetricsJSON) error
}

type metricService struct {
	logger              logging2.Logger
	storage             Storage
	metricValidator     MetricValidator
	metricLoaderService MetricLoaderService
	cfg                 *config.Config
}

func NewMetricService(logger logging2.Logger, storage Storage, metricValidator MetricValidator, metricLoaderService MetricLoaderService, cfg *config.Config) *metricService {
	ms := &metricService{
		logger:              logger,
		storage:             storage,
		metricValidator:     metricValidator,
		metricLoaderService: metricLoaderService,
		cfg:                 cfg,
	}

	isSyncMetricWrite := ms.getIsSyncMetricWrite()
	if !isSyncMetricWrite {
		go func() {
			fmt.Println("Запись на диск запущена асинхронная")
			metricLoaderService.StartMetricsToDiskInterval(context.Background(), cfg.StoreInterval)
		}()
	} else {

		fmt.Println("Запись на диск будет синхронной")
	}

	err := ms.GetMetricsFromDisk()
	if err != nil {
		logger.Error(err)
		return nil
	}

	return ms
}

func (u metricService) getIsSyncMetricWrite() bool {
	return u.cfg.StoreInterval == 0
}

func (u metricService) GetMetricsFromDisk() error {
	u.logger.Info("GetMetricsFromDisk")
	disk, err := u.metricLoaderService.LoadMetricsFromDisk(context.Background())
	if err != nil {
		u.logger.Error(err)
		return nil
	}

	u.logger.Info("GetMetricsFromDisk: Загружены данные ", disk)
	for _, k := range disk {
		u.logger.Info("GetMetricsFromDisk: Обработка элемента ", k)
		t := u.metricValidator.GetValidType(k.MType)

		switch t {
		case entity2.Gauge:
			err := u.storage.SaveGauge(context.Background(), converter.MetricJSONToMetricGauge(k))
			if err != nil {
				u.logger.Error(err)
				return err
			}
		case entity2.Counter:
			err := u.storage.SaveCounter(context.Background(), converter.MetricJSONToMetricCounter(k))
			if err != nil {
				u.logger.Error(err)
				return err
			}
		default:
			u.logger.Error("Неизвестный тип", k)
		}

	}

	return nil
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

	isSyncMetricWrite := u.getIsSyncMetricWrite()
	if isSyncMetricWrite {
		err := u.metricLoaderService.SaveMetricToDisk(ctx, converter.MetricGaugeToMetricJSON(metricsGauge))
		if err != nil {
			u.logger.Error(err)
			return err

		}
	}

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

	isSyncMetricWrite := u.getIsSyncMetricWrite()
	if isSyncMetricWrite {
		err := u.metricLoaderService.SaveMetricToDisk(ctx, converter.MetricCounterToMetricJSON(metricsCounter))
		if err != nil {
			u.logger.Error(err)
			return err

		}
	}

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
