package main

import (
	"encoding/json"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
)

func main() {
	cfg := config.NewConfig()
	cfg.Load()

	logging.Init()
	logger := logging.NewMyLogger().WithField("prefix", "SERVER")

	cfgMarshal, err := json.MarshalIndent(&cfg, "", "   ")
	if err == nil {
		logger.Info("CFG: ", string(cfgMarshal))
	}

	logger.Info("Starting metric app")
	_, err = server.New(cfg, logger)
	if err != nil {
		logger.Error("Starting metric app Failed", err)
		return
	}
}
