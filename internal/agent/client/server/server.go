package server

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"io"
	"net/http"
)

type httpClient struct {
	NewRequest    func(method, url string, body io.Reader) (*http.Request, error)
	DefaultClient http.Client
}

type client struct {
	host       string
	httpClient httpClient
}

type Client interface {
	Post(ctx context.Context, updateDto dto.Update) error
}

func New(host string) Client {
	httpClient := httpClient{
		NewRequest:    http.NewRequest,
		DefaultClient: *http.DefaultClient,
	}

	return &client{
		host:       host,
		httpClient: httpClient,
	}
}
