package service

import (
	clientServer "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

type (
	service struct {
		server     clientServer.Client
		logger     logging.Logger
		repository server.Repository
	}
)

func New(
	server clientServer.Client,
	logger logging.Logger,
	repository server.Repository,
) server.Service {
	return &service{
		server:     server,
		logger:     logger,
		repository: repository,
	}
}
