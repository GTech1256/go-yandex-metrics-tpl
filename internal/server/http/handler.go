package http

import (
	"context"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Service interface {
	GetMetrics(ctx context.Context) (string, error)
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

	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(metrics))
	if err != nil {
		h.logger.Error(err)
		return
	}
}
