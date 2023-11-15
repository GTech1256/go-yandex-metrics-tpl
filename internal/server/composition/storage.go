package composition

import (
	"context"
	"errors"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	sql2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/db/sql"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/ping"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/memory"
	sql3 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/sql"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
)

type storage interface {
	SaveGauge(ctx context.Context, gauge *entity2.MetricGauge) error
	SaveCounter(ctx context.Context, counter *entity2.MetricCounter) error
	GetGaugeValue(name string) (*entity2.GaugeValue, error)
	GetCounterValue(name string) (*entity2.CounterValue, error)
	GetAllMetrics() *metric.AllMetrics
	Ping(ctx context.Context) error
	SaveMetricBatch(ctx context.Context, metrics []*entity2.MetricJSON) error
}

var (
	ErrNoSQLConnection = errors.New("Нет подключения к БД")
)

func NewStorageComposition(cfg *config.Config, logger logging.Logger, router *chi.Mux) (storage, error) {
	var storage storage

	if cfg.GetIsEnabledSQLStore() {
		logger.Info("SQL Enabled")
		sql, err := sql2.NewSQL(*cfg.DatabaseDSN)
		//defer sql.DB.Close()
		if err != nil {
			logger.Error(err)
			return nil, ErrNoSQLConnection
		}

		storage = sql3.NewStorage(sql.DB)

		logger.Info("Регистрация /ping Router")
		pingHandler := ping.NewHandler(logger, storage)
		pingHandler.Register(router)
	} else {
		storage = memory.NewStorage()
	}

	return storage, nil
}
