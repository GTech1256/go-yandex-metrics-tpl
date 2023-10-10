package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

func main() {
	parseFlags()
	logging.Init()
	logger := logging.NewMyLogger().WithField("prefix", "SERVER")

	logger.Info("Starting server app")
	_, err := server.New(port, logger)
	logger.Infof("Started server app on %v", port)
	if err != nil {
		logger.Error("Starting server app Failed", err)
		return
	}
}
