package repository

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
)

type repository struct {
}

func New() server.Repository {
	return &repository{}
}
