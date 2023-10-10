package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/config"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

func main() {
	cfg := config.NewConfig().(*config.Config)
	cfg.Load()
	logging.Init()
	logger := logging.NewMyLogger().WithField("prefix", "AGENT")

	logger.Info("Starting agent app")
	_, err := agent.New(cfg, logger)
	if err != nil {
		logger.Error("Starting agent app Failed", err)
		return
	}
}
