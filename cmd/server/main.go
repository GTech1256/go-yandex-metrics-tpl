package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

func main() {
	cfg := config.NewConfig().(*config.Config)
	cfg.Load()
	logging.Init()
	logger := logging.NewMyLogger().WithField("prefix", "SERVER")

	logger.Info("Starting server app")
	_, err := server.New(cfg, logger)
	if err != nil {
		logger.Error("Starting server app Failed", err)
		return
	}
}
