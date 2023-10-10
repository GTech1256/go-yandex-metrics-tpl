package agent

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server"
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

func New(port string, pollInterval int, reportInterval int, logger logging2.Logger) (App, error) {
	pollIntervalDuration := time.Duration(pollInterval) * time.Second
	reportIntervalDuration := time.Duration(reportInterval) * time.Second

	router := http.NewServeMux()
	ctx := context.Background()

	metricSendCh := make(chan server2.MetricSendCh)

	serverHost := "http://localhost:8080"
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

	logger.Info(fmt.Sprintf("Start Listen Port %v", port))
	log.Fatal(http.ListenAndServe(port, logging.WithLogging(router, logger)))

	return &app{}, nil
}
