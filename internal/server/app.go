package server

import (
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/logging"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/counter"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/gauge"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"log"
	"net/http"
)

type App interface {
}

type app struct {
}

func New(port string, logger logging2.Logger) (App, error) {
	router := http.NewServeMux()

	metricStorage := metric.NewStorage()
	updateService := service.NewUpdateService(logger, metricStorage)

	logger.Info("Register /update/counter/ Router")
	updateCounterHandler := counter.NewHandler(logger, updateService)
	updateCounterHandler.Register(router)

	logger.Info("Register /update/gauge/ Router")
	updateGaugeHandler := gauge.NewHandler(logger, updateService)
	updateGaugeHandler.Register(router)

	logger.Info("Register /update/ Router")
	updateHandler := update.NewHandler(logger, updateService)
	updateHandler.Register(router)

	logger.Info(fmt.Sprintf("Start Listen Port %v", port))
	log.Fatal(http.ListenAndServe(port, logging.WithLogging(router, logger)))

	return &app{}, nil
}
