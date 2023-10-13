package service

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

func (m *MockRepository) GetMetrics() (*agentEntity.Metric, error) {
	args := m.Called()
	return args.Get(0).(*agentEntity.Metric), args.Error(1)
}
