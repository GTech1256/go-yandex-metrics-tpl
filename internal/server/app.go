package server

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	home "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/gzip"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/logging"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/ping"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/rest/counter"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/rest/gauge"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/value"
	file2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/file"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/memory"
	sql3 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/sql"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_loader"
	metricValidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	sql2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/sql"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type App struct {
}

type MetricLoaderService interface {
	StartMetricsToDiskInterval(ctx context.Context, interval time.Duration)
	LoadMetricsFromDisk(ctx context.Context) ([]*file2.MetricJSON, error)
	SaveMetricToDisk(ctx context.Context, mj *file2.MetricJSON) error
}

func New(cfg *config.Config, logger logging2.Logger) (*App, error) {
	router := chi.NewRouter()

	router.Use(gzip.WithGzip)

	metricStorage := memory.NewStorage()

	if cfg.GetIsEnabledSQLStore() {
		logger.Info("SQL Enabled")
		sql, err := sql2.NewSQL(*cfg.DatabaseDSN)
		//defer sql.DB.Close()
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		sqlStorage := sql3.NewStorage(sql.DB)

		logger.Info("Регистрация /ping Router", sql)
		pingHandler := ping.NewHandler(logger, sqlStorage)
		pingHandler.Register(router)

	}

	var metricLoaderService MetricLoaderService = nil
	// пустое значение отключает функцию записи на диск
	if cfg.GetIsEnabledFileWrite() {
		fileStorage, err := file2.NewFileStorage(*cfg.FileStoragePath)
		if err != nil {
			logger.Error("Ошибка инцилизации fileStorage ", err)
			panic(err)
		}

		metricLoaderService = metricloader.NewMetricLoaderService(logger, fileStorage, metricStorage)
	}

	validator := metricValidator.New()
	updateService := metric2.NewMetricService(logger, metricStorage, validator, metricLoaderService, cfg)

	logger.Info("Регистрация /update/counter/ Router")
	updateCounterHandler := counter.NewHandler(logger, updateService, validator)
	updateCounterHandler.Register(router)

	logger.Info("Регистрация /update/gauge/ Router")
	updateGaugeHandler := gauge.NewHandler(logger, updateService, validator)
	updateGaugeHandler.Register(router)

	logger.Info("Регистрация /update/* Router")
	updateHandler := update.NewHandler(logger, updateService, validator)
	updateHandler.Register(router)

	logger.Info("Регистрация /value/ Router")
	valueHandler := value.NewHandler(logger, updateService, validator)
	valueHandler.Register(router)

	logger.Info("Регистрация / Router")
	homeHandler := home.NewHandler(logger, updateService)
	homeHandler.Register(router)

	logger.Info(fmt.Sprintf("Start Listen Port %v", *cfg.Port))
	log.Fatal(http.ListenAndServe(*cfg.Port, logging.WithLogging(router, logger)))

	return &App{}, nil
}
