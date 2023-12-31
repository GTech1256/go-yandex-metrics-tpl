package metric

import (
	"context"
	"errors"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	dto2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/metric/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/retry"
	"time"
)

var (
	ErrSend = errors.New("метрика не отправлена")
)

func (s *service) StartReport(ctx context.Context, reportInterval time.Duration) error {
	s.logger.Info("Запуск Report")

	ticker := time.NewTicker(reportInterval)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Остановка Pool")
			ticker.Stop()
			return nil
		case <-ticker.C:
			s.logger.Info("Тик Report")

			metrics, err := s.repository.GetMetrics()
			if err != nil {
				return err
			}

			if metrics != nil {
				err := s.sendMetricRetry(ctx, metrics)
				if err != nil {
					s.logger.Error(err)
				}

			} else {
				s.logger.Info("Нет метрики для отправки")
			}

		}
	}
}

func (s *service) sendMetricRetry(ctx context.Context, metrics *entity.Metrics) error {
	err := retry.MakeRetry(
		func() error {
			err := s.sendMetric(ctx, metrics)

			if err != nil {
				s.logger.Errorf("Ошибка отправки метрики Err: %v", err)
			}

			if errors.Is(err, api.ErrRequestDo) || errors.Is(err, api.ErrInvalidResponseStatus) {
				s.logger.Errorf("Еще попытка отправить метрику")
				return err
			}

			return nil
		},
	)

	if err != nil {
		s.logger.Error(ErrSend)
		return ErrSend
	}

	s.logger.Info("Метрика отправлена")

	return nil
}

func (s *service) sendMetric(ctx context.Context, metrics *entity.Metrics) error {
	return s.sendMetricBatchController(ctx, metrics)
}

func (s *service) sendMetricBatchController(ctx context.Context, metrics *entity.Metrics) error {
	if *s.cfg.Batch {
		s.logger.Info("Отправка метрики Батчем")
		err := s.sendMetricBatch(ctx, metrics)
		if err != nil {
			s.logger.Errorf("Ошибка отправки батчем %w", err)
			return err
		}
	} else {
		if s.cfg.RateLimit != nil {
			s.sendMetricWorkerPool(ctx, metrics, *s.cfg.RateLimit)
		} else {
			for _, m := range *metrics {
				s.logger.Info("Отправка метрики")
				err := s.sendMetricItem(ctx, &m)
				if err != nil {
					s.logger.Errorf("Ошибка отправки метрики %w", err)
				}
			}
		}
	}

	return nil
}

func (s *service) sendMetricWorkerPool(ctx context.Context, metrics *entity.Metrics, workerCount int) {
	s.metricWorkerOnce.Do(func() {
		// воркер, который достает метрику из канала и отправляет
		worker := func(idx int, metricsForSendCh chan *entity.MetricFields) {
			// Обработка каждой метрики воркером
			for metric := range metricsForSendCh {
				err := s.sendMetricItem(ctx, metric)
				if err != nil {
					s.logger.Errorf("Ошибка отправки метрик по одной", err)
				}
			}
		}

		// Создание воркеров
		for w := 1; w <= workerCount; w++ {
			go worker(w, s.metricsForSendCh)
		}
	})

	// Наполнение канала метик для последующей обработки их с помощью пула воркеров
	for _, m := range *metrics {
		s.metricsForSendCh <- &m
	}
}

func (s *service) sendMetricItem(ctx context.Context, metric *entity.MetricFields) error {
	s.logger.Infof("Отправка %v", metric.MetricName)

	if err := s.server.SendUpdate(ctx, dto2.MetricFromService(metric)); err != nil {
		s.logger.Error(err)

		return err
	}

	return nil
}

func (s *service) sendMetricBatch(ctx context.Context, metrics *entity.Metrics) error {
	updateDTOs := make([]*dto.Update, 0, len(*metrics))

	for _, m := range *metrics {
		updateDTO := dto2.MetricFromService(&m)
		updateDTOs = append(updateDTOs, &updateDTO)
	}

	s.logger.Info("Отправка sendMetricBatch")
	if err := s.server.SendUpdates(ctx, updateDTOs); err != nil {
		s.logger.Error(err)

		return err
	}
	s.logger.Infof("sendMetricBatch успешно отправлена")

	return nil
}
