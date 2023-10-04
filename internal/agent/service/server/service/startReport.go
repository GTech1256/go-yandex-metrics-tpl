package service

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"time"
)

func (s *service) StartReport(ctx context.Context, metricSendCh <-chan *agentEntity.Metric, reportInterval time.Duration) error {
	s.logger.Info("Запуск Report")
	var metric *agentEntity.Metric

	go func() {
		for {
			metric = <-metricSendCh
			s.logger.Info("Метрика получена")
		}
	}()

	for {
		<-time.After(reportInterval)
		s.logger.Info("Тик Report")

		if metric != nil {
			for _, m := range *metric {
				s.logger.Info("Отправка метрики")
				err := s.SendMetric(ctx, &m)
				if err != nil {
					return err
				}
			}
		} else {
			s.logger.Info("Нет метрики для отправки")
		}
	}
}
