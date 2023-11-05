package ping

import (
	"context"
	http2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Storage interface {
	Ping(ctx context.Context) error
}

type handler struct {
	logger  logging2.Logger
	storage Storage
}

func NewHandler(logger logging2.Logger, storage Storage) http2.Handler {
	return &handler{
		logger:  logger,
		storage: storage,
	}
}

func (h handler) Register(router *chi.Mux) {
	router.Get("/ping", h.Ping)
}

func (h handler) Ping(writer http.ResponseWriter, request *http.Request) {
	err := h.storage.Ping(context.Background())

	if err != nil {
		h.logger.Error(err)

		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}
