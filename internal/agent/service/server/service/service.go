package service

import (
	clientServer "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
	"github.com/sirupsen/logrus"
)

type (
	service struct {
		server     clientServer.Client
		logger     *logrus.Entry
		repository server.Repository
	}
)

func New(
	server clientServer.Client,
	logger *logrus.Entry,
	repository server.Repository,
) server.Service {
	return &service{
		server:     server,
		logger:     logger,
		repository: repository,
	}
}
