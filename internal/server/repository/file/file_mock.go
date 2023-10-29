package file

import "github.com/stretchr/testify/mock"

// MockFileStorage - Мок-объект для интерфейса FileStorage.
type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) ReadAll() ([]*MetricsJSON, error) {
	args := m.Called()
	return args.Get(0).([]*MetricsJSON), args.Error(1)
}

func (m *MockFileStorage) Write(metric *MetricsJSON) error {
	args := m.Called(metric)
	return args.Error(0)
}

func (m *MockFileStorage) Truncate() error {
	args := m.Called()
	return args.Error(0)
}
