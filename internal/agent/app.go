package agent

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/config"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server/repository"
	serverService "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server/service"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"time"
)

type App struct {
}

func New(cfg *config.Config, logger logging2.Logger) (*App, error) {
	pollIntervalDuration := time.Duration(*cfg.PollInterval) * time.Second
	reportIntervalDuration := time.Duration(*cfg.ReportInterval) * time.Second

	ctx := context.Background()

	serverHost := fmt.Sprintf("http://%v", *cfg.ServerPort)
	serverClient := server.New(serverHost, logger)

	serverRepository := repository.New()
	service := serverService.New(serverClient, logger, serverRepository)

	go func() {
		err := service.StartPoll(ctx, pollIntervalDuration)
		if err != nil {
			logger.Error("Ошибка начала сбора метрик", err)
		}
	}()

	go func() {
		err := service.StartReport(ctx, reportIntervalDuration)
		if err != nil {
			logger.Error("Ошибка начала отправки метрик", err)
		}
	}()

	return &App{}, nil
}
