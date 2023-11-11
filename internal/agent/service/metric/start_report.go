package metric

import (
	"context"
	"errors"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	dto2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/metric/dto"
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
				//err := s.sendMetricBatch(ctx, metrics)
				//if err != nil {
				//	return err
				//}
				for _, m := range *metrics {
					s.logger.Info("Отправка метрики")
					err := s.sendMetric(ctx, &m)
					if err != nil {
						s.logger.Errorf("Ошибка отпраки метрики %v", err)
					}
				}
			} else {
				s.logger.Info("Нет метрики для отправки")
			}

		}
	}
}

func (s *service) sendMetric(ctx context.Context, metric *entity.MetricFields) error {
	s.logger.Infof("Отправка %v", metric.MetricName)

	if err := s.server.SendUpdate(ctx, dto2.MetricFromService(metric)); err != nil {
		s.logger.Errorf("Ошибка отправки %v %v", metric.MetricName, err)

		return ErrSend
	}

	return nil
}

func (s *service) sendMetricBatch(ctx context.Context, metric *entity.Metric) error {
	updateDTOs := make([]*dto.Update, 0, len(*metric))

	for _, m := range *metric {
		updateDTO := dto2.MetricFromService(&m)
		updateDTOs = append(updateDTOs, &updateDTO)
	}

	s.logger.Infof("Отправка sendMetricBatch")
	if err := s.server.SendUpdates(ctx, updateDTOs); err != nil {
		s.logger.Errorf("Ошибка отправки")

		return ErrSend
	}

	return nil
}
