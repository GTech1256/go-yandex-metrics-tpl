package update

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service"
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
	router.HandleFunc("/update/", h.Update)
}

// Update /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (h handler) Update(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.logger.Info("Allow Only Method POST ")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.updateService.GetMetric(context.Background(), request.RequestURI)

	if err == service.ErrNotCorrectValue {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err == service.ErrNotCorrectType {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if metric.Type == entity.Gauge {
		metricsGauge := &entity.MetricGauge{
			Type:  metric.Type,
			Name:  metric.Name,
			Value: entity.GaugeValue(*metric.Value.(*float64)),
		}

		if err := h.updateService.SaveGaugeMetric(context.Background(), metricsGauge); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if metric.Type == entity.Counter {
		metricsCounter := &entity.MetricCounter{
			Type:  metric.Type,
			Name:  metric.Name,
			Value: entity.CounterValue(*metric.Value.(*int64)),
		}

		if err := h.updateService.SaveCounterMetric(context.Background(), metricsCounter); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	writer.Header().Add("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}
