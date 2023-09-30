package update

import (
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http/update/middlware/guard"
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
	router.Handle("/update/", guard.WithMetricGuarding(http.HandlerFunc(h.Update), h.logger))
}

// Update /update/
// /update/counter/ и /update/gauge обрабатываются во внутренних хендлерах
// Из-за этого остается только StatusBadRequest отдавать в остальных случаях
func (h handler) Update(writer http.ResponseWriter, request *http.Request) {
	h.logger.Info("/update/ -> http.StatusBadRequest")
	writer.WriteHeader(http.StatusBadRequest)
}
