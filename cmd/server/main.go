package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server"
	logger2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

const (
	PORT = ":8080"
)

func main() {
	logger2.Init()
	logger := logger2.GetLogger().WithField("prefix", "SERVER")
	logger.Info()
	logger.Info("Starting app")
	_, err := internal.New(PORT, logger)
	logger.Info("Started app on ")
	if err != nil {
		logger.Info("Starting app Failed", err)
		return
	}
}
