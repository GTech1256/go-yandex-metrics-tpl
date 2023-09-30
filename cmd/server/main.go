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

	//logger.SetOutput(ioutil.Discard)
	logger.Info()
	logger.Info("Starting app")
	_, err := internal.New(PORT, logger)
	logger.Info("Started app on ")
	if err != nil {
		logger.Info("Starting app Failed", err)
		return
	}
}
