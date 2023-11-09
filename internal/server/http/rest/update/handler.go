package update

import (
	"context"
	"encoding/json"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/middlware/guard"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric/converter"
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
	SaveMetricJSON(ctx context.Context, metric *entity.MetricJSON) error
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
	router.Post("/update/", h.Update)
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
	ctx := request.Context()

	// Получение метрики из запроса
	var m *entity.MetricJSON
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&m)
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Сохранение метрики
	err = h.service.SaveMetricJSON(ctx, m)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Формирование JSON ответа
	res, err := h.getMetricResponse(ctx, *m)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ответ
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		h.logger.Error(err)
		return
	}
}

func (h handler) getMetricResponse(ctx context.Context, m entity.MetricJSON) ([]byte, error) {
	// Получение метрики из сервиса
	v := converter.MetricJSONToMetricValueDTO(m)
	newData, err := h.service.GetMetricValue(ctx, &v)
	if err != nil {
		return nil, err
	}

	if m.MType == "counter" {
		int, err := strconv.ParseInt(*newData, 10, 64)
		if err != nil {
			return nil, err
		}
		m.Delta = &int
	} else if m.MType == "gauge" {
		float, err := strconv.ParseFloat(*newData, 64)
		if err != nil {
			return nil, err
		}
		m.Value = &float
	}

	return json.Marshal(newData)
}
