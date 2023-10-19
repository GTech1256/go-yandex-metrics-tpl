package update

import (
	"context"
	"encoding/json"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/converter"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/middlware/guard"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/models"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type MetricValidator interface {
	GetValidType(metricType string) entity2.Type
	MakeMetricValuesFromURL(url string) (*entity2.MetricFields, error)
}

type Service interface {
	SaveCounterMetric(ctx context.Context, metric *entity.MetricFields) error
	SaveGaugeMetric(ctx context.Context, metric *entity.MetricFields) error
	GetMetricValue(ctx context.Context, metric *updateInterface.GetMetricValueDto) (*string, error)
}

type handler struct {
	logger          logging2.Logger
	metricValidator MetricValidator
	service         Service
}

func NewHandler(logger logging2.Logger, service Service, metricValidator MetricValidator) http2.Handler {
	return &handler{
		logger:          logger,
		service:         service,
		metricValidator: metricValidator,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Post("/update", h.Update)
	router.Handle("/update/*", guard.WithMetricGuarding(http.HandlerFunc(h.UpdateRest), h.logger, h.metricValidator))
}

// UpdateRest /update/
// /update/counter/ и /update/gauge обрабатываются во внутренних хендлерах
// Из-за этого остается только StatusBadRequest отдавать в остальных случаях
func (h handler) UpdateRest(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
}

// Update POST /update
func (h handler) Update(writer http.ResponseWriter, request *http.Request) {
	var m *models.Metrics
	ctx := context.Background()

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&m)
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mType := h.metricValidator.GetValidType(m.MType)

	switch mType {
	case entity2.Gauge:
		mg := converter.MetricsGaugeToMetricFields(*m)
		err := h.service.SaveGaugeMetric(ctx, &mg)
		if err != nil {
			h.logger.Error(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		mv := converter.MetricsToMetricValueDTO(*m)
		value, err := h.service.GetMetricValue(ctx, &mv)
		if err != nil {
			h.logger.Error(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		valueFloat, err := strconv.ParseFloat(*value, 64)
		if err != nil {
			h.logger.Error(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.Value = &valueFloat

		break
	case entity2.Counter:
		mc := converter.MetricsCounterToMetricFields(*m)
		err := h.service.SaveCounterMetric(ctx, &mc)
		if err != nil {
			h.logger.Error(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		mv := converter.MetricsToMetricValueDTO(*m)
		value, err := h.service.GetMetricValue(ctx, &mv)
		if err != nil {
			h.logger.Error(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		valueInt, err := strconv.ParseInt(*value, 10, 64)
		if err != nil {
			h.logger.Error(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.Delta = &valueInt
		break
	default:
		h.logger.Error("Неизвестный тип метрики ", m)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(m)
	if err != nil {
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		h.logger.Error(err)
		return
	}
}
