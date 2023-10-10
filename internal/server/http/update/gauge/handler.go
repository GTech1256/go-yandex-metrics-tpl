package gauge

import (
	"context"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/middlware/guard"
	metricvalidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type handler struct {
	logger          logging2.Logger
	updateService   updateInterface.Service
	metricValidator metricvalidator.MetricValidator
}

func NewHandler(logger logging2.Logger, updateService updateInterface.Service, metricValidator metricvalidator.MetricValidator) http2.Handler {
	return &handler{
		logger:          logger.WithField("TYPE", "HANDLER").WithField("METRIC", "GAUGE"),
		updateService:   updateService,
		metricValidator: metricValidator,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Post("/update/gauge/{name}/{value}", guard.WithMetricGuarding(http.HandlerFunc(h.UpdateGauge), h.logger, h.metricValidator))
}

// UpdateGauge /update/gauge/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (h handler) UpdateGauge(writer http.ResponseWriter, request *http.Request) {
	metricFields, err := h.metricValidator.MakeMetricValuesFromURL(request.RequestURI)
	if err != nil {
		h.logger.Error("При получении полей метрик gauge из URL произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.updateService.SaveGaugeMetric(context.Background(), metricFields); err != nil {
		h.logger.Error("При сохранении метрики произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Add("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}
