package value

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type Service interface {
	GetMetricValue(ctx context.Context, metric *updateInterface.GetMetricValueDto) (*string, error)
}

type handler struct {
	logger        logging2.Logger
	updateService Service
}

func NewHandler(logger logging2.Logger, updateService Service) http2.Handler {
	return &handler{
		logger:        logger,
		updateService: updateService,
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

	trimmed := *value

	if metric.Type == string(entity.Gauge) {
		trimmed = strings.TrimRight(strings.TrimRight(*value, "0"), ".")
	}

	_, err = writer.Write([]byte(trimmed))
	if err != nil {
		h.logger.Error(err)
	}
	writer.WriteHeader(http.StatusOK)
}
