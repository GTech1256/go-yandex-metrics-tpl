package sql

import (
	"context"
	"errors"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/jackc/pgx/v5"
)

func (s *storage) SaveCounter(ctx context.Context, counter *entity2.MetricCounter) error {
	return s.saveCounter(ctx, counter, s.db)
}

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s *storage) saveCounter(ctx context.Context, counter *entity2.MetricCounter, executor Executor) error {
	oldValue, err := s.GetCounterValue(counter.Name)
	isNoOldValue := errors.Is(err, pgx.ErrNoRows)

	if err != nil && !isNoOldValue {
		return err
	}

	if isNoOldValue {
		err = s.insertCounter(ctx, counter, executor)
		if err != nil {
			return err
		}

		return nil
	}

	c := &entity2.MetricCounter{
		Type:  counter.Type,
		Name:  counter.Name,
		Value: *oldValue + counter.Value,
	}

	err = s.updateCounter(ctx, c, executor)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) insertCounter(ctx context.Context, counter *entity2.MetricCounter, executor Executor) error {
	_, err := executor.Exec(ctx, InsertCounter, counter.Name, counter.Value)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) updateCounter(ctx context.Context, counter *entity2.MetricCounter, executor Executor) error {
	_, err := executor.Exec(ctx, UpdateCounter, counter.Name, counter.Value)
	if err != nil {
		return err
	}

	return nil
}
