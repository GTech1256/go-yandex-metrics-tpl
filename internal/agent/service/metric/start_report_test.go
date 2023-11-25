package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/config"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	mock2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/service/metric/mock"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_service_StartReport(t *testing.T) {
	t.Skipf("skip")
	// Arrange
	ctx := context.Background()
	reportInterval := time.Duration(5) * time.Millisecond
	mockMetric := &agentEntity.Metrics{
		agentEntity.MetricFields{
			MetricType:  "testType",
			MetricName:  "testName",
			MetricValue: "testValue",
		},
	}

	repo := new(mock2.MockRepository)
	client := new(mock2.MockClient)
	mockLogger := new(logging.LoggerMock)
	cfg := config.NewConfig()

	s := New(client, mockLogger, repo, cfg)

	// Assert
	repo.On("GetMetrics").Return(mockMetric, nil)

	// TODO: Проверять s.On("sendMetric")
	// сейчас в s.sendMetric вызывается s.client.Post из-за этого есть эта проверка
	client.On("SendUpdate", ctx, mock.MatchedBy(func(update dto.Update) bool {
		return update.Type == (*mockMetric)[0].MetricType &&
			update.Name == (*mockMetric)[0].MetricName &&
			update.Value == (*mockMetric)[0].MetricValue
	})).Return(nil)

	mockLogger.On("Info", "Запуск Report")
	mockLogger.On("Info", "Тик Report")
	mockLogger.On("Info", "Отправка метрики")
	mockLogger.On("Info", "Отправка sendMetricBatch")
	mockLogger.On("Infof", []interface{}{"Отправка %v", (*mockMetric)[0].MetricName})

	// Act
	errCh := make(chan error)
	go func() {
		errCh <- s.StartReport(ctx, reportInterval)
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
	repo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
