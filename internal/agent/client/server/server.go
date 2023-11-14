package server

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	serverHttp "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api/update"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api/updates"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

type UpdateAPI interface {
	SendUpdate(ctx context.Context, updateDto dto.Update) error
	SendUpdateJSON(ctx context.Context, updateDto dto.Update) error
}

type UpdatesAPI interface {
	SendUpdates(ctx context.Context, updateDto []*dto.Update) error
}

type client struct {
	host       string
	httpClient serverHttp.ClientHTTP
	logger     logging.Logger
	updateAPI  UpdateAPI
	updatesAPI UpdatesAPI
}

func New(host string, logger logging.Logger) *client {
	httpClient := serverHttp.New()
	updateAPI := update.New(httpClient, host, logger)
	updatesAPI := updates.New(httpClient, host, logger)

	return &client{
		host:       host,
		httpClient: httpClient,
		logger:     logger,
		updateAPI:  updateAPI,
		updatesAPI: updatesAPI,
	}
}

func (c *client) SendUpdate(ctx context.Context, updateDto dto.Update) error {
	return c.updateAPI.SendUpdate(ctx, updateDto)
}

func (c *client) SendUpdateJSON(ctx context.Context, updateDto dto.Update) error {
	return c.updateAPI.SendUpdateJSON(ctx, updateDto)
}

func (c *client) SendUpdates(ctx context.Context, updateDto []*dto.Update) error {
	return c.updatesAPI.SendUpdates(ctx, updateDto)
}
