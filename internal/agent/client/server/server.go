package server

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

type client struct {
	host       string
	httpClient HTTPClient
	logger     logging.Logger
}

type Client interface {
	Post(ctx context.Context, updateDto dto.Update) error
}

func New(host string, logger logging.Logger) Client {
	httpClient := &httpClient{}

	return &client{
		host:       host,
		httpClient: httpClient,
		logger:     logger,
	}
}
