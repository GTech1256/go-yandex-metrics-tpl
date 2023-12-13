package updates

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/client"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

type updates struct {
	HTTPClient client.ClientHTTP
	BaseURL    string
	logger     logging.Logger
}

func New(HTTPClient client.ClientHTTP, BaseURL string, logger logging.Logger) *updates {
	return &updates{
		HTTPClient: HTTPClient,
		BaseURL:    BaseURL,
		logger:     logger,
	}
}
