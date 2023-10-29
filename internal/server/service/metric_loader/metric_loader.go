package metricloader

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/file"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"time"
)

type FileStorage interface {
	ReadAll() ([]*file.MetricsJSON, error)
	Write(metric *file.MetricsJSON) error
	Truncate() error
}

type MetricProvider interface {
	GetAllMetrics() *metric2.AllMetrics
}
type metricLoaderService struct {
	logger         logging2.Logger
	storage        FileStorage
	metricProvider MetricProvider
}

func NewMetricLoaderService(logger logging2.Logger, storage FileStorage, metricProvider MetricProvider) *metricLoaderService {
	return &metricLoaderService{
		logger:         logger,
		storage:        storage,
		metricProvider: metricProvider,
	}
}

func (u metricLoaderService) LoadMetricsFromDisk(ctx context.Context) ([]*file.MetricsJSON, error) {
	all, err := u.storage.ReadAll()
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return all, nil
}

func (u metricLoaderService) StartMetricsToDiskInterval(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			u.logger.Info("Остановка Интервала для записи метрик на диск")
			ticker.Stop()
		case <-ticker.C:
			err := u.saveMetricsToDisk(ctx)
			if err != nil {
				u.logger.Error(err)
				return
			}
		}
	}
}

func (u metricLoaderService) saveMetricsToDisk(ctx context.Context) error {
	metrics := u.metricProvider.GetAllMetrics()
	u.logger.Info("Сохранение всех Метрик на диск ", metrics)

	err := u.storage.Truncate()
	if err != nil {
		u.logger.Error("Ошибка при очистке файла на диске ", err)
		return err
	}

	for k, v := range metrics.Counter {
		Delta := v
		metricJSON := &file.MetricsJSON{
			ID:    k,
			MType: string(entity.Counter),
			Delta: &Delta,
		}

		err := u.storage.Write(metricJSON)
		if err != nil {
			u.logger.Error(err)
			return err
		}
	}

	for k, v := range metrics.Gauge {
		Value := v
		metricJSON := &file.MetricsJSON{
			ID:    k,
			MType: string(entity.Gauge),
			Value: &Value,
		}

		err := u.storage.Write(metricJSON)
		if err != nil {
			u.logger.Error(err)
			return err
		}
	}

	u.logger.Info("Все Метрики сохранены на диск ", metrics)

	return nil
}

func (u metricLoaderService) SaveMetricToDisk(ctx context.Context, mj *file.MetricsJSON) error {
	u.logger.Info("Сохранение Метрики на диск ", mj)
	err := u.storage.Write(mj)
	if err != nil {
		u.logger.Error("Метрика не сохранена на диск ", err, mj)
		return err
	}
	u.logger.Info("Метрика сохранена на диск ", mj)

	return nil
}

func (u metricLoaderService) clear() error {
	u.logger.Info("Отчистка fileStorage")
	return u.storage.Truncate()
}
