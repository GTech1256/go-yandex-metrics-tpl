package gauge

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http/update/middlware/guard"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

type handler struct {
	logger        *logrus.Entry
	updateService updateInterface.Service
}

func NewHandler(logger *logrus.Entry, updateService updateInterface.Service) http2.Handler {
	return &handler{
		logger:        logger,
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
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metricCounterValue, err := util.GetTypeGaugeValue(metricFields.MetricValue)
	// При попытке передать запрос с некорректным типом метрики
	// или значением возвращать http.StatusBadRequest.
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metricsCounter := &entity.MetricGauge{
		Type:  entity.Counter,
		Name:  metricFields.MetricName,
		Value: entity.GaugeValue(*metricCounterValue),
	}

	if err := h.updateService.SaveGaugeMetric(context.Background(), metricsCounter); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Add("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}
