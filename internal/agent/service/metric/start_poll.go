package metric

import (
	"context"
	"time"
)

func (s *service) StartPoll(ctx context.Context, pollInterval time.Duration) error {
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

			err := s.repository.LoadMetric(ctx)
			if err != nil {
				return err
			}
		}
	}
}
