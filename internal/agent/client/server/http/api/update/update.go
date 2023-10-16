package update

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

type update struct {
	HTTPClient http.ClientHTTP
	BaseURL    string
	logger     logging.Logger
}

func New(HTTPClient http.ClientHTTP, BaseURL string, logger logging.Logger) *update {
	return &update{
		HTTPClient: HTTPClient,
		BaseURL:    BaseURL,
		logger:     logger,
	}
}
