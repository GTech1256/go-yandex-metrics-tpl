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
	fmt.Println("SQL")
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	start := time.Now()
	for _, metric := range metrics {
		if metric.GetIsGauge() {
			s.logger.Info("Save metric Gauge %+v", metric)
			err := s.saveGauge(ctx, converter.MetricJSONToMetricGauge(metric), tx)
			if err != nil {
				s.logger.Error("Save metric Gauge %+v", metric)
				return err
			}
		} else if metric.GetIsCounter() {
			s.logger.Info("Save metric Counter %+v", metric)
			err := s.saveCounter(ctx, converter.MetricJSONToMetricCounter(metric), tx)
			if err != nil {
				s.logger.Error("Save metric Counter %+v", metric)
				return err
			}
		} else {
			s.logger.Infof("Unknown metric %+v", metric)
		}
	}
	fmt.Println("SaveMetricBatch:", time.Since(start))

	err = retry.MakeRetry(
		func() error {
			err = tx.Commit(ctx)

			// Ошибка чтения данных из сети или БД из-за проблем соединения.
			if err != nil {
				s.logger.Errorf("Ошибка применении транзакции %v", err)
				return err
			}

			return nil
		},
	)

	//err = tx.Commit(ctx)

	if err != nil {

		fmt.Println(33)
		s.logger.Errorf("Не удалось сохранить метрику батчем %v", err)
		return err
	}

	s.logger.Info("Метрик сохранена батчем")

	return nil
}
