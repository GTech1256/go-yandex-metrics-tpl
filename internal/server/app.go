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
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service"
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

	metricStorage := metric.NewStorage()
	validator := metricValidator.New()
	updateService := service.NewUpdateService(logger, metricStorage, validator)

	logger.Info("Register /update/counter/ Router")
	updateCounterHandler := counter.NewHandler(logger, updateService, validator)
	updateCounterHandler.Register(router)

	logger.Info("Register /update/gauge/ Router")
	updateGaugeHandler := gauge.NewHandler(logger, updateService, validator)
	updateGaugeHandler.Register(router)

	logger.Info("Register /update/* Router")
	updateHandler := update.NewHandler(logger, updateService, validator)
	updateHandler.Register(router)

	logger.Info("Register /value/ Router")
	valueHandler := value.NewHandler(logger, updateService, validator)
	valueHandler.Register(router)

	logger.Info("Register / Router")
	homeHandler := home.NewHandler(logger, updateService)
	homeHandler.Register(router)

	logger.Info(fmt.Sprintf("Start Listen Port %v", *cfg.Port))
	log.Fatal(http.ListenAndServe(*cfg.Port, logging.WithLogging(router, logger)))

	return &App{}, nil
}
