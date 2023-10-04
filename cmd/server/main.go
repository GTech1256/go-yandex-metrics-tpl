package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

const (
	PORT = ":8080"
)

func main() {
	logging.Init()
	logger := logging.GetLogger().WithField("prefix", "SERVER")

	logger.Info("Starting server app")
	_, err := server.New(PORT, logger)
	logger.Info("Started server app on ")
	if err != nil {
		logger.Error("Starting server app Failed", err)
		return
	}
}
