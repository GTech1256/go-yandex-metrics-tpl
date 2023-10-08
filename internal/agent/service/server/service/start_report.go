package service

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
	"time"
)

//var (
//	ErrReceiveEmptyMetric = errors.New("получена пустая метрика")
//)

func (s *service) StartReport(ctx context.Context, metricSendCh <-chan server.MetricSendCh, reportInterval time.Duration) error {
	s.logger.Info("Запуск Report")

	ticker := time.NewTicker(reportInterval)
	var metric *agentEntity.Metric

	go func() {
		for {
			data, ok := <-metricSendCh
			if !ok {
				ctx.Done()
				break
			}

			metric = data.Data

			if metric == nil {
				s.logger.Errorf("Получена пустая метрика")
			}
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
