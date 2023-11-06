package sql

import (
	"context"
	"database/sql"
	"fmt"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/jackc/pgx/v5"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type storage struct {
	db *pgx.Conn
}

func NewStorage(db *pgx.Conn) *storage {

	return &storage{
		db: db,
	}
}

// SaveGauge новое значение должно замещать предыдущее.
func (s *storage) SaveGauge(ctx context.Context, gauge *entity2.MetricGauge) error {
	_, err := s.GetGaugeValue(gauge.Name)
	hasGauge := !(err == sql.ErrNoRows)
	if err != nil && !hasGauge {
		return err
	}

	if hasGauge {
		err := s.updateGauge(ctx, gauge)
		if err != nil {
			return err
		}

		return nil
	}

	err = s.insertGauge(ctx, gauge)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) insertGauge(ctx context.Context, gauge *entity2.MetricGauge) error {
	query := "INSERT INTO gauge(title, value) VALUES($1, $2)"
	_, err := s.db.Exec(ctx, query, gauge.Name, gauge.Value)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) updateGauge(ctx context.Context, gauge *entity2.MetricGauge) error {
	query := "UPDATE gauge SET value = $2 where title = $1"
	_, err := s.db.Exec(ctx, query, gauge.Name, gauge.Value)
	if err != nil {
		return err
	}

	return nil
}

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s *storage) SaveCounter(ctx context.Context, counter *entity2.MetricCounter) error {
	oldValue, err := s.GetCounterValue(counter.Name)
	hasCounter := !(err == sql.ErrNoRows)
	if err != nil && !hasCounter {
		return err
	}

	if hasCounter {
		c := &entity2.MetricCounter{
			Type:  counter.Type,
			Name:  counter.Name,
			Value: *oldValue + counter.Value,
		}

		err = s.updateCounter(ctx, c)
		if err != nil {
			return err
		}

		return nil
	}

	err = s.insertCounter(ctx, counter)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) insertCounter(ctx context.Context, counter *entity2.MetricCounter) error {
	query := "INSERT INTO counter(title, delta) VALUES($1, $2)"
	_, err := s.db.Exec(ctx, query, counter.Name, counter.Value)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) updateCounter(ctx context.Context, counter *entity2.MetricCounter) error {
	query := "UPDATE counter SET delta = $2 where title = $1"
	_, err := s.db.Exec(ctx, query, counter.Name, counter.Value)
	if err != nil {
		return err
	}

	return nil
}

// GetGauge - возвращает значение Gauge из хранилища
func (s *storage) GetGaugeValue(name string) (*entity2.GaugeValue, error) {
	ctx := context.Background()

	var v *entity2.GaugeValue
	query := "SELECT value FROM gauge where title = $1"
	row := s.db.QueryRow(ctx, query, name)

	err := row.Scan(&v)
	if err != nil {
		return nil, err
	}

	fmt.Println(v)

	return v, nil
}

// GetCounter - возвращает значение Counter из хранилища
func (s *storage) GetCounterValue(name string) (*entity2.CounterValue, error) {
	ctx := context.Background()

	var v *entity2.CounterValue
	query := "SELECT value FROM counter where title = $1"
	row := s.db.QueryRow(ctx, query, name)

	err := row.Scan(&v)
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
