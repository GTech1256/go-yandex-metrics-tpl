package updates

import (
	"context"
	"encoding/json"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type MetricValidator interface {
	GetValidType(metricType string) entity.Type
	MakeMetricValuesFromURL(url string) (*entity.MetricFields, error)
}

type Service interface {
	SaveMetricJSONs(ctx context.Context, metrics []*entity.MetricJSON) error
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
	ctx := request.Context()

	// Получение метрики из запроса
	var m []*entity.MetricJSON

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&m)
	if err != nil {
		h.logger.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Сохранение метрики
	err = h.service.SaveMetricJSONs(ctx, m)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Ответ
	writer.WriteHeader(http.StatusOK)
	if err != nil {
		h.logger.Error(err)
		return
	}
}
