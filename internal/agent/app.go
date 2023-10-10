package agent

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/config"
	server2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server/repository"
	serverService "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server/service"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/logging"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"log"
	"net/http"
	"time"
)

type App interface {
}

type app struct {
}

func New(cfg *config.Config, logger logging2.Logger) (App, error) {
	pollIntervalDuration := time.Duration(*cfg.PollInterval) * time.Second
	reportIntervalDuration := time.Duration(*cfg.ReportInterval) * time.Second

	router := http.NewServeMux()
	ctx := context.Background()

	metricSendCh := make(chan server2.MetricSendCh)

	serverHost := fmt.Sprintf("http://+%v", *cfg.ServerPort)
	serverClient := server.New(serverHost, logger)
	serverRepository := repository.New()
	service := serverService.New(serverClient, logger, serverRepository)

	go func() {
		err := service.StartPoll(ctx, metricSendCh, pollIntervalDuration)
		if err != nil {
			logger.Error("Ошибка начала сборка метрик", err)
		}
	}()

	go func() {
		err := service.StartReport(ctx, metricSendCh, reportIntervalDuration)
		if err != nil {
			logger.Error("Ошибка начала отправки метрик", err)
		}
	}()

	logger.Info(fmt.Sprintf("Start Listen Port %v", *cfg.AgentPort))
	log.Fatal(http.ListenAndServe(*cfg.AgentPort, logging.WithLogging(router, logger)))

	return &app{}, nil
}
