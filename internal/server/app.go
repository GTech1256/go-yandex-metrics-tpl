package server

import (
	"fmt"
	home "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/logging"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/counter"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/gauge"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/value"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service"
	metricvalidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type App interface {
}

type app struct {
}

func New(port string, logger logging2.Logger) (App, error) {
	router := chi.NewRouter()

	metricStorage := metric.NewStorage()
	metricValidator := metricvalidator.New()
	updateService := service.NewUpdateService(logger, metricStorage, metricValidator)

	logger.Info("Register /update/counter/ Router")
	updateCounterHandler := counter.NewHandler(logger, updateService, metricValidator)
	updateCounterHandler.Register(router)

	logger.Info("Register /update/gauge/ Router")
	updateGaugeHandler := gauge.NewHandler(logger, updateService, metricValidator)
	updateGaugeHandler.Register(router)

	logger.Info("Register /update/* Router")
	updateHandler := update.NewHandler(logger, updateService, metricValidator)
	updateHandler.Register(router)

	logger.Info("Register /value/ Router")
	valueHandler := value.NewHandler(logger, updateService, metricValidator)
	valueHandler.Register(router)

	logger.Info("Register / Router")
	homeHandler := home.NewHandler(logger, updateService, metricValidator)
	homeHandler.Register(router)

	logger.Info(fmt.Sprintf("Start Listen Port %v", port))
	log.Fatal(http.ListenAndServe(port, logging.WithLogging(router, logger)))

	return &app{}, nil
}
