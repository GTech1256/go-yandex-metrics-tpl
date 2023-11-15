package mock

import (
	"context"
	agentEntity "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"github.com/stretchr/testify/mock"
)

// Mocked Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) LoadMetric(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRepository) GetMetrics() (*agentEntity.Metrics, error) {
	args := m.Called()
	return args.Get(0).(*agentEntity.Metrics), args.Error(1)
}
