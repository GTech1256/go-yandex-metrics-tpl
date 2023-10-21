package value

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/models"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/value/converter"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
)

type MetricValidator interface {
	GetValidType(metricType string) entity.Type
}

type Service interface {
	GetMetricValue(ctx context.Context, metric *updateInterface.GetMetricValueDto) (*string, error)
}

type handler struct {
	logger          logging2.Logger
	updateService   Service
	metricValidator MetricValidator
}

func NewHandler(logger logging2.Logger, updateService Service, metricValidator MetricValidator) http2.Handler {
	return &handler{
		logger:          logger,
		updateService:   updateService,
		metricValidator: metricValidator,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Get("/value/{type}/{name}", h.Value)
	router.Post("/value/", h.ValueJSON)
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

func (h handler) ValueJSON(writer http.ResponseWriter, request *http.Request) {
	var m *models.Metrics
	ctx := context.Background()

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&m)
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mv := converter.MetricsToMetricValueDTO(*m)
	value, err := h.updateService.GetMetricValue(ctx, &mv)
	if err != nil {
		h.logger.Error(err, 1)
	}
	if value == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	mType := h.metricValidator.GetValidType(m.MType)

	switch mType {
	case entity.Gauge:
		valueFloat, err := strconv.ParseFloat(*value, 64)
		if err != nil {
			h.logger.Error(err, 2)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.Value = &valueFloat
	case entity.Counter:
		valueInt, err := strconv.ParseInt(*value, 10, 64)
		if err != nil {
			h.logger.Error(err, 3)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.Delta = &valueInt
	default:
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("HELLO", m)

	res, err := json.Marshal(m)
	if err != nil {
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		h.logger.Error(err)
	}
}
