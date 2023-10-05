package gauge

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/middlware/guard"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/util"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"net/http"
)

type handler struct {
	logger        logging2.Logger
	updateService updateInterface.Service
}

func NewHandler(logger logging2.Logger, updateService updateInterface.Service) http2.Handler {
	return &handler{
		logger:        logger.WithField("TYPE", "HANDLER").WithField("METRIC", "GAUGE"),
		updateService: updateService,
	}
}

func (h handler) Register(router *http.ServeMux) {
	router.Handle("/update/gauge/", guard.WithMetricGuarding(http.HandlerFunc(h.UpdateGauge), h.logger))
}

// UpdateGauge /update/gauge/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (h handler) UpdateGauge(writer http.ResponseWriter, request *http.Request) {
	metricFields, err := util.MakeMetricValuesFromURL(request.RequestURI)
	if err != nil {
		h.logger.Error("При получении полей метрик gauge из URL произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: to Service
	metricCounterValue, err := util.GetTypeGaugeValue(metricFields.MetricValue)
	// При попытке передать запрос с некорректным типом метрики
	// или значением возвращать service.StatusBadRequest.
	if err != nil {
		h.logger.Error("При получении значения метрики произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metricsGauge := &entity.MetricGauge{
		Type:  entity.Counter,
		Name:  metricFields.MetricName,
		Value: entity.GaugeValue(*metricCounterValue),
	}

	if err := h.updateService.SaveGaugeMetric(context.Background(), metricsGauge); err != nil {
		h.logger.Error("При сохранении метрики произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Error("Метрика сохранена", metricsGauge)
	writer.Header().Add("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}
