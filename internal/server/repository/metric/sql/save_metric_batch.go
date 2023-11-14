package sql

import (
	"context"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/sql/converter"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s *storage) SaveMetricBatch(ctx context.Context, metrics []*entity2.MetricJSON) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		if metric.GetIsGauge() {
			err := s.saveGauge(ctx, converter.MetricJSONToMetricGauge(metric), tx)
			if err != nil {
				return err
			}
		} else if metric.GetIsCounter() {
			err := s.saveCounter(ctx, converter.MetricJSONToMetricCounter(metric), tx)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
