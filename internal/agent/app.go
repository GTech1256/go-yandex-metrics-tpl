package agent

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server/repository"
	serverService "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server/service"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/logging"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

type App interface {
}

type app struct {
}

const (
	pollInterval   time.Duration = time.Duration(2) * time.Second
	reportInterval time.Duration = time.Duration(10) * time.Second
)

func New(port string, logger *logrus.Entry) (App, error) {
	router := http.NewServeMux()

	ctx := context.Background()

	metricSendCh := make(chan *agentEntity.Metric)

	serverHost := "http://localhost:8080"
	serverClient := server.New(serverHost)
	serverRepository := repository.New()
	service := serverService.New(serverClient, logger, serverRepository)

	go func() {
		err := service.StartPoll(ctx, metricSendCh, pollInterval)
		if err != nil {
			logger.Error("Ошибка начала сборка метрик", err)
		}
	}()

	go func() {
		err := service.StartReport(ctx, metricSendCh, reportInterval)
		if err != nil {
			logger.Error("Ошибка начала отправки метрик", err)
		}
	}()

	logger.Info(fmt.Sprintf("Start Listen Port %v", port))
	log.Fatal(http.ListenAndServe(port, logging.WithLogging(router, logger)))

	return &app{}, nil
}
