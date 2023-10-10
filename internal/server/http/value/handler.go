package value

import (
	"context"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type handler struct {
	logger          logging2.Logger
	updateService   updateInterface.Service
	metricValidator metricValidator.MetricValidator
}

func NewHandler(logger logging2.Logger, updateService updateInterface.Service, metricValidator metricValidator.MetricValidator) http2.Handler {
	return &handler{
		logger:          logger,
		updateService:   updateService,
		metricValidator: metricValidator,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Get("/value/{type}/{name}", h.Value)
}

// Value /value/{type}/{name}
func (h handler) Value(writer http.ResponseWriter, request *http.Request) {
	metric := &updateInterface.GetMetricValueDto{
		Type: chi.URLParam(request, "type"),
		Name: chi.URLParam(request, "name"),
	}

	value, err := h.updateService.GetMetricValue(context.Background(), metric)
	if err != nil {
		h.logger.Error(err)
	}
	if value == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	trimmed := strings.TrimRight(strings.TrimRight(*value, "0"), ".")

	_, err = writer.Write([]byte(trimmed))
	if err != nil {
		h.logger.Error(err)
	}
	writer.WriteHeader(http.StatusOK)
}
