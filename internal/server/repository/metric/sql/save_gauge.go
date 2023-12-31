package sql

import (
	"context"
	"errors"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/jackc/pgx/v5"
)

func (s *storage) SaveGauge(ctx context.Context, gauge *entity2.MetricGauge) error {
	return s.saveGauge(ctx, gauge, s.db)
}

// SaveGauge новое значение должно замещать предыдущее.
func (s *storage) saveGauge(ctx context.Context, gauge *entity2.MetricGauge, executor Executor) error {
	_, err := s.getGaugeValue(gauge.Name, executor)
	isUniqueValue := errors.Is(err, pgx.ErrNoRows)
	if err != nil && !isUniqueValue {
		return err
	}

	if isUniqueValue {
		err = s.insertGauge(ctx, gauge, executor)
		if err != nil {
			return err
		}

		return nil
	}

	err = s.updateGauge(ctx, gauge, executor)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) insertGauge(ctx context.Context, gauge *entity2.MetricGauge, executor Executor) error {
	_, err := executor.Exec(ctx, InsertGauge, gauge.Name, gauge.Value)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) updateGauge(ctx context.Context, gauge *entity2.MetricGauge, executor Executor) error {
	_, err := executor.Exec(ctx, UpdateGauge, gauge.Name, gauge.Value)
	if err != nil {
		return err
	}

	return nil
}
