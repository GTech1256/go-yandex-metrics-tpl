package service

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/stretchr/testify/mock"
)

// Mocked Client
type MockClient struct {
	mock.Mock
	PostParam dto.Update
}

func (m *MockClient) Post(ctx context.Context, updateDto dto.Update) error {
	m.PostParam = updateDto

	args := m.Called(ctx, updateDto)

	return args.Error(0)
}
