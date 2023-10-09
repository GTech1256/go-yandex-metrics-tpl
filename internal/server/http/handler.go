package http

import (
	"context"

	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type handler struct {
	logger          logging2.Logger
	updateService   updateInterface.Service
	metricValidator metric_validator.MetricValidator
}

func NewHandler(logger logging2.Logger, updateService updateInterface.Service, metricValidator metric_validator.MetricValidator) Handler {
	return &handler{
		logger:          logger,
		updateService:   updateService,
		metricValidator: metricValidator,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Get("/", h.Home)
}

// Home /
func (h handler) Home(writer http.ResponseWriter, request *http.Request) {
	metrics, err := h.updateService.GetMetrics(context.Background())
	if err != nil {
		h.logger.Error(err)
		return
	}

	_, err = writer.Write([]byte(metrics))
	if err != nil {
		h.logger.Error(err)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
