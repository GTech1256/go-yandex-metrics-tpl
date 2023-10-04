package service

import (
	"context"
	"errors"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"time"
)

var (
	ErrSend = errors.New("метрика не отправлена")
)

func (s *service) StartPoll(ctx context.Context, metricSendCh chan<- *agentEntity.Metric, pollInterval time.Duration) error {
	s.logger.Info("Запуск Pool")
	for {
		<-time.After(pollInterval)
		s.logger.Info("Тик Pool")
		metric, err := s.repository.GetMetric(ctx)

		if err != nil {
			return err
		}
		s.logger.Info("Отправка agent.Metric")
		metricSendCh <- metric

	}
}

func (s *service) SendMetric(ctx context.Context, metric *entity.MetricFields) error {
	s.logger.Info("Отправка ", metric.MetricName)

	if err := s.server.Post(ctx, dto.Update{
		Type:  metric.MetricType,
		Name:  metric.MetricName,
		Value: metric.MetricValue,
	}); err != nil {
		return ErrSend
	}

	return nil
}
