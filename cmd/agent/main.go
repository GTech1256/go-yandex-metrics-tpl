package main

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
)

func main() {
	parseFlags()
	logging.Init()
	logger := logging.NewMyLogger().WithField("prefix", "AGENT")

	logger.Info("Starting agent app")
	_, err := agent.New(port, pollInterval, reportInterval, logger)
	logger.Infof("Started agent app on %v", port)
	if err != nil {
		logger.Error("Starting agent app Failed", err)
		return
	}
}
