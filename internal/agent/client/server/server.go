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
	updateApi  UpdateAPI
	updatesApi UpdatesAPI
}

func New(host string, logger logging.Logger) *client {
	httpClient := serverHttp.New()
	updateApi := update.New(httpClient, host, logger)
	updatesApi := updates.New(httpClient, host, logger)

	return &client{
		host:       host,
		httpClient: httpClient,
		logger:     logger,
		updateApi:  updateApi,
		updatesApi: updatesApi,
	}
}

func (c *client) SendUpdate(ctx context.Context, updateDto dto.Update) error {
	return c.updateApi.SendUpdate(ctx, updateDto)
}

func (c *client) SendUpdateJSON(ctx context.Context, updateDto dto.Update) error {
	return c.updateApi.SendUpdateJSON(ctx, updateDto)
}

func (c *client) SendUpdates(ctx context.Context, updateDto []*dto.Update) error {
	return c.updatesApi.SendUpdates(ctx, updateDto)
}
