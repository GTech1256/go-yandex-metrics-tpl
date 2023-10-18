package server

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	serverHttp "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api/update"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

type UpdateAPI interface {
	SendUpdate(ctx context.Context, updateDto dto.Update) error
}

type client struct {
	host       string
	httpClient serverHttp.ClientHTTP
	logger     logging.Logger
	api        UpdateAPI
}

func New(host string, logger logging.Logger) *client {
	httpClient := serverHttp.New()
	api := update.New(httpClient, host, logger)

	return &client{
		host:       host,
		httpClient: httpClient,
		logger:     logger,
		api:        api,
	}
}

func (c *client) SendUpdate(ctx context.Context, updateDto dto.Update) error {
	return c.api.SendUpdate(ctx, updateDto)
}
