package memory

import (
	"context"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/stretchr/testify/mock"
)

// MockMetricProvider - Мок-объект для интерфейса MetricProvider.
type MockMetricProvider struct {
	mock.Mock
}

func (m *MockMetricProvider) GetAllMetrics(ctx context.Context) *metric2.AllMetrics {
	args := m.Called()
	return args.Get(0).(*metric2.AllMetrics)
}
