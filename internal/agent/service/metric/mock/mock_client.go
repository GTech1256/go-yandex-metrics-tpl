package mock

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/stretchr/testify/mock"
)

// Mocked Client
type MockClient struct {
	mock.Mock
	PostParam  dto.Update
	PostParams []*dto.Update
}

func (m *MockClient) SendUpdate(ctx context.Context, updateDto dto.Update) error {
	m.PostParam = updateDto

	args := m.Called(ctx, updateDto)

	return args.Error(0)
}

func (m *MockClient) SendUpdateJSON(ctx context.Context, updateDto dto.Update) error {
	m.PostParam = updateDto

	args := m.Called(ctx, updateDto)

	return args.Error(0)
}

func (m *MockClient) SendUpdates(ctx context.Context, updateDto []*dto.Update) error {
	m.PostParams = updateDto

	args := m.Called(ctx, updateDto)

	return args.Error(0)
}
