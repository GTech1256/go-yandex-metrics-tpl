package update

import (
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/middlware/guard"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"net/http"
)

type handler struct {
	logger        logging2.Logger
	updateService updateInterface.Service
}

func NewHandler(logger logging2.Logger, updateService updateInterface.Service) http2.Handler {
	return &handler{
		logger:        logger,
		updateService: updateService,
	}
}

func (h handler) Register(router *http.ServeMux) {
	router.Handle("/update/", guard.WithMetricGuarding(http.HandlerFunc(h.Update), h.logger))
}

// Update /update/
// /update/counter/ и /update/gauge обрабатываются во внутренних хендлерах
// Из-за этого остается только StatusBadRequest отдавать в остальных случаях
func (h handler) Update(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
}
