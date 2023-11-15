package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	mock2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/metric/mock"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_service_StartPoll(t *testing.T) {
	// Arrange
	ctx := context.Background()
	updateInterval := time.Millisecond * 100

	repo := new(mock2.MockRepository)

	client := new(mock2.MockClient)

	mockLogger := new(logging.LoggerMock)

	// Assert
	repo.On("LoadMetric", ctx).Return(nil)
	mockLogger.On("Info", "Запуск Pool")
	mockLogger.On("Info", "Тик Pool")

	service := New(client, mockLogger, repo)

	// Act
	go func() {
		err := service.StartPoll(ctx, updateInterval)
		assert.NoError(t, err)
	}()

	// Ожидание прогона service.StartPoll
	<-time.After(updateInterval * 2)

	// Clean up
	ctx.Done()
	repo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)

}

func Test_service_sendMetricItem(t *testing.T) {
	ctx := context.Background()
	// создаем тестовый контекст и метрику
	testMetric := &entity.MetricFields{
		MetricType:  "testType",
		MetricName:  "testName",
		MetricValue: "testValue",
	}

	repo := new(mock2.MockRepository)

	client := new(mock2.MockClient)
	client.On("SendUpdate", ctx, mock.MatchedBy(func(update dto.Update) bool {
		return update.Type == "testType" &&
			update.Name == "testName" &&
			update.Value == "testValue"
	})).Return(nil)

	mockLogger := new(logging.LoggerMock)
	mockLogger.On("Infof", []interface{}{"Отправка %v", testMetric.MetricName})

	s := &service{client, mockLogger, repo}

	// вызываем функцию sendMetric
	err := s.sendMetricItem(ctx, testMetric)

	// проверяем, что ошибки нет
	assert.NoError(t, err)

	// Clean up
	repo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
