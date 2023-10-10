package service

import (
	"context"
	"errors"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	server2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"time"
)

var (
	ErrSend = errors.New("метрика не отправлена")
)

func (s *service) StartPoll(ctx context.Context, metricSendCh chan<- server2.MetricSendCh, pollInterval time.Duration) error {
	s.logger.Info("Запуск Pool")

	ticker := time.NewTicker(pollInterval)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Остановка Pool")
			ticker.Stop()
			return nil

		case <-ticker.C:

			s.logger.Info("Тик Pool")
			metric, err := s.repository.GetMetric(ctx)

			if err != nil {
				return err
			}
			s.logger.Info("Отправка Pool метрики в канал")
			metricSendCh <- server2.MetricSendCh{
				ID:   "StartPoll fn",
				Data: metric,
			}
		}
	}
}

func (s *service) sendMetric(ctx context.Context, metric *entity.MetricFields) error {
	s.logger.Infof("Отправка %v", metric.MetricName)

	if err := s.server.Post(ctx, dto.Update{
		Type:  metric.MetricType,
		Name:  metric.MetricName,
		Value: metric.MetricValue,
	}); err != nil {
		s.logger.Errorf("Ошибка отправки %v", metric.MetricName)

		return ErrSend
	}

	return nil
}
