package server

import (
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	home "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/gzip"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/logging"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/rest/counter"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/rest/gauge"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/value"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_loader"
	metricValidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type App struct {
}

func New(cfg *config.Config, logger logging2.Logger) (*App, error) {
	router := chi.NewRouter()

	router.Use(gzip.WithGzip)

	fileStorage, err := metric.NewFileStorage(*cfg.FileStoragePath)
	if err != nil {
		logger.Error("Ошибка инцилизации fileStorage ", err)
		panic(err)
	}

	metricStorage := metric.NewStorage()

	metricLoaderService := metric_loader.NewMetricLoaderService(logger, fileStorage, metricStorage)

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
