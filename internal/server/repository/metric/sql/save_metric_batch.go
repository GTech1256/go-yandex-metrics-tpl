package sql

import (
	"context"
	"fmt"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/sql/converter"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/retry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

var (
	ErrMetricJSONToMetricGaugeConvertation   = "Конвертировать MetricJSON в MetricGauge не удалось"
	ErrMetricJSONToMetricCounterConvertation = "Конвертировать MetricJSON в MetricCounter не удалось"
)

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s *storage) SaveMetricBatch(ctx context.Context, metrics []*entity2.MetricJSON) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	start := time.Now()
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
	fmt.Println("SaveMetricBatch:", time.Since(start))

	err = retry.MakeRetry(
		func() error {
			err = tx.Commit(ctx)

			// Ошибка чтения данных из сети или БД из-за проблем соединения.
			if err != nil {
				s.logger.Errorf("Ошибка применении транзакции %w", err)
				return err
			}

			return nil
		},
	)

	if err != nil {
		s.logger.Errorf("Не удалось сохранить метрику батчем %w", err)
		return err
	}

	s.logger.Errorf("Удалось сохранить метрику батчем %w", err)
	return nil
}
