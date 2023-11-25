package mock

import (
	"github.com/stretchr/testify/mock"
)

// MockConfig is a mock implementation of the Config interface
type MockConfig struct {
	mock.Mock
}

// Load mocks the Load method of the Config interface
func (m *MockConfig) Load() {
	m.Called()
}

// GetServerPort mocks the GetServerPort method of the Config interface
func (m *MockConfig) GetServerPort() *string {
	args := m.Called()
	return args.Get(0).(*string)
}

// GetReportInterval mocks the GetReportInterval method of the Config interface
func (m *MockConfig) GetReportInterval() *int {
	args := m.Called()
	return args.Get(0).(*int)
}

// GetPollInterval mocks the GetPollInterval method of the Config interface
func (m *MockConfig) GetPollInterval() *int {
	args := m.Called()
	return args.Get(0).(*int)
}

// GetBatch mocks the GetBatch method of the Config interface
func (m *MockConfig) GetBatch() *bool {
	args := m.Called()
	return args.Get(0).(*bool)
}
