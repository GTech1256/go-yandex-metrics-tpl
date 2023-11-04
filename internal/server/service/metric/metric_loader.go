package metric

import (
	"context"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/file"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric/converter"
)

func (u metricService) getIsSyncMetricWrite() bool {
	return u.cfg.StoreInterval == 0
}

func (u metricService) GetMetricsFromDisk() error {
	u.logger.Info("GetMetricsFromDisk")
	metricsFromDisk, err := u.metricLoaderService.LoadMetricsFromDisk(context.Background())
	if err != nil {
		u.logger.Error(err)
		return nil
	}

	u.logger.Info("GetMetricsFromDisk: Загружены данные ", metricsFromDisk)
	for _, metricJSON := range metricsFromDisk {
		u.logger.Info("GetMetricsFromDisk: Обработка элемента ", metricJSON)
		metricType := u.metricValidator.GetValidType(metricJSON.MType)

		switch metricType {
		case entity2.Gauge:
			err := u.storage.SaveGauge(context.Background(), converter.MetricJSONToMetricGauge(metricJSON))
			if err != nil {
				u.logger.Error(err)
				return err
			}
		case entity2.Counter:
			err := u.storage.SaveCounter(context.Background(), converter.MetricJSONToMetricCounter(metricJSON))
			if err != nil {
				u.logger.Error(err)
				return err
			}
		default:
			u.logger.Error("Неизвестный тип", metricJSON)
		}

	}

	return nil
}

func (u metricService) saveMetricCallback(ctx context.Context, metricJSON *file.MetricJSON) {
	if u.metricLoaderService == nil {
		return
	}

	isSyncMetricWrite := u.getIsSyncMetricWrite()
	if isSyncMetricWrite {
		u.logger.Info("Метрика сохраняется в синхронной стратегии")
		err := u.metricLoaderService.SaveMetricToDisk(ctx, metricJSON)
		if err != nil {
			u.logger.Error(err)
		}
	}
}

func (u metricService) initMetricLoader() {
	if u.metricLoaderService == nil {
		return
	}

	isSyncMetricWrite := u.getIsSyncMetricWrite()
	if !isSyncMetricWrite {
		go func() {
			u.logger.Info("Запись на диск запущена асинхронная")
			u.metricLoaderService.StartMetricsToDiskInterval(context.Background(), u.cfg.StoreInterval)
		}()
	} else {
		u.logger.Info("Запись на диск будет синхронной")
	}

	err := u.GetMetricsFromDisk()
	if err != nil {
		u.logger.Error(err)
	}

}
