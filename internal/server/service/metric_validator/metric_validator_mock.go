package metricValidator

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

var _ MetricValidator = (*MockMetricValidator)(nil)

// Mocked MetricValidator
type MockMetricValidator struct {
	mock.Mock
}

func (m *MockMetricValidator) MakeMetricValuesFromURL(url string) (*entity.MetricFields, error) {
	args := m.Called(url)

	return args.Get(0).(*entity.MetricFields), args.Error(1)
}

func (m *MockMetricValidator) GetValidType(metricType string) entity.Type {
	args := m.Called(metricType)

	return args.Get(0).(entity.Type)
}

func (m *MockMetricValidator) GetTypeGaugeValue(metricValueUnsafe string) (*float64, error) {
	args := m.Called(metricValueUnsafe)

	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockMetricValidator) GetTypeCounterValue(metricValueUnsafe string) (*int64, error) {
	args := m.Called(metricValueUnsafe)

	return args.Get(0).(*int64), args.Error(1)
}
