package http

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type Service interface {
	GetMetrics(ctx context.Context) (*metric.AllMetrics, error)
}

type handler struct {
	logger        logging2.Logger
	updateService Service
}

func NewHandler(logger logging2.Logger, updateService Service) Handler {
	return &handler{
		logger:        logger,
		updateService: updateService,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Get("/", h.Home)
}

// Home /
func (h handler) Home(writer http.ResponseWriter, request *http.Request) {
	metrics, err := h.updateService.GetMetrics(context.Background())
	if err != nil {
		h.logger.Error(err)
		return
	}

	gaugeMetrics := make([]string, len(metrics.Gauge))
	counterMetrics := make([]string, len(metrics.Counter))

	for name, value := range metrics.Gauge {
		gaugeMetrics = append(gaugeMetrics, fmt.Sprintf("<tr>"+
			"<td>%v</td>"+
			"<td>%v</td>"+
			"</tr>", name, value),
		)
	}
	for name, value := range metrics.Counter {
		counterMetrics = append(counterMetrics, fmt.Sprintf("<tr>"+
			"<td>%v</td>"+
			"<td>%v</td>"+
			"</tr>", name, value),
		)
	}

	metricsList := fmt.Sprintf("<h1>MetricsJSON</h1>"+
		"<div style='display: flex;width: 800px;justify-content: space-between;'>"+
		"<div><h2>Gauge</h2><table><tr><th>Name</th><th>Value</th></tr>%v</table></div>"+
		"<div><h2>Counter</h2><table><tr><th>Name</th><th>Value</th></tr>%v</table></div>"+
		"</div>", strings.Join(gaugeMetrics, ""), strings.Join(counterMetrics, ""))

	html := fmt.Sprintf("%v", metricsList)

	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(html))
	if err != nil {
		h.logger.Error(err)
		return
	}
}
