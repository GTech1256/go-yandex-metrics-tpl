package update

import (
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/middlware/guard"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type MetricValidator interface {
	MakeMetricValuesFromURL(url string) (*entity2.MetricFields, error)
}

type handler struct {
	logger          logging2.Logger
	metricValidator MetricValidator
}

func NewHandler(logger logging2.Logger, metricValidator MetricValidator) http2.Handler {
	return &handler{
		logger:          logger,
		metricValidator: metricValidator,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Handle("/update/*", guard.WithMetricGuarding(http.HandlerFunc(h.Update), h.logger, h.metricValidator))
}

// Update /update/
// /update/counter/ и /update/gauge обрабатываются во внутренних хендлерах
// Из-за этого остается только StatusBadRequest отдавать в остальных случаях
func (h handler) Update(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
}
