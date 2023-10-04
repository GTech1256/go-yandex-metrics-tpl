package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

const (
	PORT = ":8081"
)

func main() {
	logging.Init()
	logger := logging.GetLogger().WithField("prefix", "AGENT")

	logger.Info("Starting agent app")
	_, err := agent.New(PORT, logger)
	logger.Info("Started agent app on ")
	if err != nil {
		logger.Error("Starting agent app Failed", err)
		return
	}
}
