package service

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"time"
)

func (s *service) StartReport(ctx context.Context, metricSendCh <-chan *agentEntity.Metric, reportInterval time.Duration) error {
	s.logger.Info("Запуск Report")

	ticker := time.NewTicker(reportInterval)
	var metric *agentEntity.Metric

	go func() {
		for {
			metric = <-metricSendCh
			s.logger.Info("Метрика получена")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Остановка Pool")
			ticker.Stop()
			return nil
		case <-ticker.C:
			s.logger.Info("Тик Report")

			if metric != nil {
				for _, m := range *metric {
					s.logger.Info("Отправка метрики")
					err := s.sendMetric(ctx, &m)
					if err != nil {
						return err
					}
				}
			} else {
				s.logger.Info("Нет метрики для отправки")
			}

		}
	}
}
