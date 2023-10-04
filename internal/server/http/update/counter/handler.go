package counter

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/middlware/guard"
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
		logger:        logger.WithField("TYPE", "HANDLER").WithField("METRIC", "COUNTER"),
		updateService: updateService,
	}
}

func (h handler) Register(router *http.ServeMux) {
	router.Handle("/update/counter/", guard.WithMetricGuarding(http.HandlerFunc(h.UpdateCounter), h.logger))
}

// UpdateCounter /update/counter/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (h handler) UpdateCounter(writer http.ResponseWriter, request *http.Request) {
	metricFields, err := util.MakeMetricValuesFromURL(request.RequestURI)
	if err != nil {
		h.logger.Error("При получении полей метрик из URL произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metricCounterValue, err := util.GetTypeCounterValue(metricFields.MetricValue)
	// При попытке передать запрос с некорректным типом метрики
	// или значением возвращать service.StatusBadRequest.
	if err != nil {
		h.logger.Error("При получении значения метрики произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metricsCounter := &entity.MetricCounter{
		Type:  entity.Counter,
		Name:  metricFields.MetricName,
		Value: entity.CounterValue(*metricCounterValue),
	}

	if err := h.updateService.SaveCounterMetric(context.Background(), metricsCounter); err != nil {
		h.logger.Error("При сохранении метрики произошла ошибка ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	h.logger.Error("Метрика сохранена", metricsCounter)
	writer.Header().Add("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}
