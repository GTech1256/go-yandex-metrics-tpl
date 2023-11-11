package updates

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

type updates struct {
	HTTPClient http.ClientHTTP
	BaseURL    string
	logger     logging.Logger
}

func New(HTTPClient http.ClientHTTP, BaseURL string, logger logging.Logger) *updates {
	return &updates{
		HTTPClient: HTTPClient,
		BaseURL:    BaseURL,
		logger:     logger,
	}
}
