package service

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
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
	mockMetric := agentEntity.Metric{
		entity.MetricFields{
			MetricType:  "testType",
			MetricName:  "testName",
			MetricValue: "testValue",
		},
	}
	metricSendCh := make(chan *agentEntity.Metric)
	reportInterval := time.Duration(2) * time.Second

	repo := new(MockRepository)

	client := new(MockClient)
	client.On("Post", ctx, mock.MatchedBy(func(update dto.Update) bool {
		return update.Type == "testType" &&
			update.Name == "testName" &&
			update.Value == "testValue"
	})).Return(nil)

	mockLogger := new(logging.LoggerMock)
	mockLogger.On("Info", "Запуск Report")
	mockLogger.On("Info", "Метрика получена")
	mockLogger.On("Info", "Тик Report")
	mockLogger.On("Info", "Отправка метрики")
	mockLogger.On("Info", "Отправка ", mockMetric[0].MetricName)

	s := New(client, mockLogger, repo)

	errCh := make(chan error)
	go func() {
		errCh <- s.StartReport(ctx, metricSendCh, reportInterval)

	}()

	go func() {
		metricSendCh <- &mockMetric
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
	close(metricSendCh)
	repo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
