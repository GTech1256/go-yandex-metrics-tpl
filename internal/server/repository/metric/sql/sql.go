package sql

import (
	"context"
	"errors"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type storage struct {
	logger logging.Logger
	db     DB
}

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Ping(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewStorage(db DB, logger logging.Logger) *storage {
	return &storage{
		db:     db,
		logger: logger,
	}
}

// GetGauge - возвращает значение Gauge из хранилища
func (s *storage) GetGaugeValue(name string) (*entity2.GaugeValue, error) {
	ctx := context.Background()

	var v *entity2.GaugeValue
	query := "SELECT value FROM gauge where title = $1"
	row := s.db.QueryRow(ctx, query, name)

	err := row.Scan(&v)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetCounter - возвращает значение Counter из хранилища
func (s *storage) GetCounterValue(name string) (*entity2.CounterValue, error) {
	ctx := context.Background()

	var v *entity2.CounterValue
	query := "SELECT delta FROM counter where title = $1"
	row := s.db.QueryRow(ctx, query, name)

	err := row.Scan(&v)

	//if errors.Is(err, pgx.ErrNoRows) {
	//	return nil, nil
	//}

	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *storage) getGaugeMetrics(ctx context.Context) (*map[string]float64, error) {
	v := make(map[string]float64, 0)

	query := "SELECT title, value FROM gauge"
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			title string
			value float64
		)

		err := rows.Scan(&title, &value)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		v[title] = value
	}

	return &v, nil
}

func (s *storage) getCounterMetrics(ctx context.Context) (*map[string]int64, error) {
	v := make(map[string]int64, 0)

	query := "SELECT title, delta FROM counter"
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			title string
			delta int64
		)

		err := rows.Scan(&title, &delta)
		if err != nil {
			return nil, err
		}

		v[title] = delta
	}

	return &v, nil
}

func (s *storage) GetAllMetrics() *metric.AllMetrics {
	ctx := context.Background()

	gauge, err := s.getGaugeMetrics(ctx)
	if err != nil {
		return nil
	}

	counter, err := s.getCounterMetrics(ctx)
	if err != nil {
		return nil
	}

	return &metric.AllMetrics{
		Gauge:   *gauge,
		Counter: *counter,
	}
}

func (s *storage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
