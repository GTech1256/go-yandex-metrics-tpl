package metric

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/stretchr/testify/mock"
)

// MockService - мок для интерфейса Service
type MockService struct {
	mock.Mock
}

// SaveGaugeMetric - мок для метода SaveGaugeMetric
func (m *MockService) SaveGaugeMetric(ctx context.Context, metric *entity.MetricFields) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

// SaveCounterMetric - мок для метода SaveCounterMetric
func (m *MockService) SaveCounterMetric(ctx context.Context, metric *entity.MetricFields) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

// GetMetricValue - мок для метода GetMetricValue
func (m *MockService) GetMetricValue(ctx context.Context, metric *updateinterface.GetMetricValueDto) (*string, error) {
	args := m.Called(ctx, metric)
	if strPtr, ok := args.Get(0).(*string); ok {
		return strPtr, args.Error(1)
	}
	return nil, args.Error(1)
}

// GetMetrics - мок для метода GetMetrics
func (m *MockService) GetMetrics(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
