package service

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	server2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_service_StartReport(t *testing.T) {
	// Arrange
	ctx := context.Background()
	reportInterval := time.Duration(5) * time.Millisecond
	mockMetric := &agentEntity.Metric{
		entity.MetricFields{
			MetricType:  "testType",
			MetricName:  "testName",
			MetricValue: "testValue",
		},
	}
	metricSendCh := make(chan server2.MetricSendCh)

	repo := new(MockRepository)

	client := new(MockClient)
	client.On("Post", ctx, mock.MatchedBy(func(update dto.Update) bool {
		return update.Type == "testType" &&
			update.Name == "testName" &&
			update.Value == "testValue"
	})).Return(nil)

	mockLogger := new(logging.LoggerMock)
	mockLogger.On("Info", "Запуск Report")
	mockLogger.On("Info", "Тик Report")
	mockLogger.On("Info", "Отправка метрики")
	mockLogger.On("Info", "Отправка ", (*mockMetric)[0].MetricName)

	s := New(client, mockLogger, repo)

	errCh := make(chan error)
	go func() {
		errCh <- s.StartReport(ctx, metricSendCh, reportInterval)
	}()

	go func() {
		fmt.Println("TEST SEND METRIC", mockMetric)
		metricSendCh <- server2.MetricSendCh{
			Id:   "Test_service_StartReport fn",
			Data: mockMetric,
		}
	}()

	// ждем, чтобы метод StartReport выполнился
	time.Sleep(reportInterval * 3)

	// проверяем, что ошибки нет
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	default:
	}

	// Clean up
	ctx.Done()
	close(metricSendCh)
	repo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
