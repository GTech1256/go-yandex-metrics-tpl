package service

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/server"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_service_StartPoll(t *testing.T) {
	// Arrange
	ctx := context.Background()
	updateInterval := time.Millisecond * 100
	mockMetric := &agentEntity.Metric{
		entity.MetricFields{
			MetricType:  "type",
			MetricName:  "name",
			MetricValue: "value",
		},
	}

	repo := new(MockRepository)
	repo.On("GetMetric", ctx).Return(mockMetric, nil)

	client := new(MockClient)

	mockLogger := new(logging.LoggerMock)
	mockLogger.On("Info", "Запуск Pool")
	mockLogger.On("Info", "Тик Pool")
	mockLogger.On("Info", "Отправка Pool метрики в канал")

	service := New(client, mockLogger, repo)

	metricSendCh := make(chan server.MetricSendCh)

	// Act
	go func() {
		err := service.StartPoll(ctx, metricSendCh, updateInterval)
		assert.NoError(t, err)
	}()

	// Assert
	select {
	case receivedMetric := <-metricSendCh:
		assert.Equal(t, mockMetric, receivedMetric.Data)
	case <-time.After(updateInterval * 2): // Give it some time to run
		t.Error("Timed out waiting for metric to be sent")
	}

	// Clean up
	ctx.Done()
	close(metricSendCh)
	repo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)

}

func Test_service_sendMetric(t *testing.T) {
	ctx := context.Background()
	// создаем тестовый контекст и метрику
	testMetric := &entity.MetricFields{
		MetricType:  "testType",
		MetricName:  "testName",
		MetricValue: "testValue",
	}

	repo := new(MockRepository)

	client := new(MockClient)
	client.On("Post", ctx, mock.MatchedBy(func(update dto.Update) bool {
		return update.Type == "testType" &&
			update.Name == "testName" &&
			update.Value == "testValue"
	})).Return(nil)

	mockLogger := new(logging.LoggerMock)
	mockLogger.On("Infof", []interface{}{"Отправка %v", testMetric.MetricName})

	s := &service{client, mockLogger, repo}

	// вызываем функцию sendMetric
	err := s.sendMetric(ctx, testMetric)

	// проверяем, что ошибки нет
	assert.NoError(t, err)

	// Clean up
	repo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
