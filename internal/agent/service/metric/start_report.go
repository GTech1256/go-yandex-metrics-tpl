package metric

import (
	"context"
	"errors"
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

	if err := s.server.SendUpdateJSON(ctx, dto2.MetricFromService(metric)); err != nil {
		s.logger.Errorf("Ошибка отправки %v", metric.MetricName)

		return ErrSend
	}

	return nil
}
