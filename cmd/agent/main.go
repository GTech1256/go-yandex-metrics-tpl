package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/app"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/config"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

func main() {
	cfg := config.NewConfig()
	cfg.Load()

	logging.Init()
	logger := logging.NewMyLogger().WithField("prefix", "AGENT")

	logger.Info("Starting agent app")
	_, err := app.New(cfg, logger)

	// TODO: Добавить Graceful Shutdown
	// Сейчас остается чтобы сервис сразу после запуска не завершался
	shutdown := make(chan int)
	<-shutdown

	if err != nil {
		logger.Error("Starting agent app Failed", err)
		return
	}
}
