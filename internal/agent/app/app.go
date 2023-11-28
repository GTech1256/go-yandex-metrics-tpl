package app

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/config"
	metricRepository "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/repository/metric"
	metricService "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"time"
)

type App struct {
}

func New(cfg *config.Config, logger logging.Logger) (*App, error) {
	pollIntervalDuration := time.Duration(*cfg.PollInterval) * time.Second
	reportIntervalDuration := time.Duration(*cfg.ReportInterval) * time.Second

	ctx := context.Background()

	serverHost := fmt.Sprintf("http://%v", *cfg.ServerPort)
	serverClient := server.New(serverHost, logger)

	mr := metricRepository.New()
	ms := metricService.New(serverClient, logger, mr, cfg)

	go func() {
		err := ms.StartPoll(ctx, pollIntervalDuration)
		if err != nil {
			logger.Error("Ошибка начала сбора метрик", err)
		}
	}()

	go func() {
		err := ms.StartReport(ctx, reportIntervalDuration)
		if err != nil {
			logger.Error("Ошибка начала отправки метрик", err)
		}
	}()

	return &App{}, nil
}
