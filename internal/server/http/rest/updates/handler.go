package updates

import (
	"context"
	"encoding/json"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/converter"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type MetricValidator interface {
	GetValidType(metricType string) entity.Type
	MakeMetricValuesFromURL(url string) (*entity.MetricFields, error)
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
	router.Post("/updates", h.Updates)
}

// Updates POST /updates
func (h handler) Updates(writer http.ResponseWriter, request *http.Request) {
	var ms *[]*entity.MetricJSON
	ctx := context.Background()

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&ms)
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, m := range *ms {

		mType := h.metricValidator.GetValidType(m.MType)

		switch mType {
		case entity.Gauge:
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
		case entity.Counter:
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
		default:
			h.logger.Error("Неизвестный тип метрики ", m)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//
	//res, err := json.Marshal(m)
	//if err != nil {
	//	return
	//}
	//
	//writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	//_, err = writer.Write()
	//if err != nil {
	//	h.logger.Error(err)
	//	return
	//}
}
