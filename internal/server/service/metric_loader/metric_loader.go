package metricloader

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/file"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"time"
)

type FileStorage interface {
	ReadAll() ([]*file.MetricJSON, error)
	Write(metric *file.MetricJSON) error
	Truncate() error
}

type MetricProvider interface {
	GetAllMetrics() *metric2.AllMetrics
}
type MetricLoaderService struct {
	logger         logging2.Logger
	storage        FileStorage
	metricProvider MetricProvider
}

func NewMetricLoaderService(logger logging2.Logger, storage FileStorage, metricProvider MetricProvider) *MetricLoaderService {
	return &MetricLoaderService{
		logger:         logger,
		storage:        storage,
		metricProvider: metricProvider,
	}
}

func (u MetricLoaderService) LoadMetricsFromDisk(ctx context.Context) ([]*file.MetricJSON, error) {
	all, err := u.storage.ReadAll()
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return all, nil
}

func (u MetricLoaderService) StartMetricsToDiskInterval(ctx context.Context, interval time.Duration) {
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

func (u MetricLoaderService) saveMetricsToDisk(ctx context.Context) error {
	metrics := u.metricProvider.GetAllMetrics()
	u.logger.Info("Сохранение всех Метрик на диск ", metrics)

	err := u.storage.Truncate()
	if err != nil {
		u.logger.Error("Ошибка при очистке файла на диске ", err)
		return err
	}

	for name, value := range metrics.Counter {
		Delta := value
		metricJSON := &file.MetricJSON{
			ID:    name,
			MType: string(entity.Counter),
			Delta: &Delta,
		}

		err := u.storage.Write(metricJSON)
		if err != nil {
			u.logger.Error(err)
			return err
		}
	}

	for name, value := range metrics.Gauge {
		Value := value
		metricJSON := &file.MetricJSON{
			ID:    name,
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

func (u MetricLoaderService) SaveMetricToDisk(ctx context.Context, mj *file.MetricJSON) error {
	u.logger.Info("Сохранение Метрики на диск ", mj)
	err := u.storage.Write(mj)
	if err != nil {
		u.logger.Error("Метрика не сохранена на диск ", err, mj)
		return err
	}
	u.logger.Info("Метрика сохранена на диск ", mj)

	return nil
}

func (u MetricLoaderService) clear() error {
	u.logger.Info("Отчистка fileStorage")
	return u.storage.Truncate()
}
