package metricvalidator

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricValidator_MakeMetricValuesFromURL(t *testing.T) {
	validator := New()

	t.Run("Допустимый URL-адрес", func(t *testing.T) {
		url := "/update/counter/metric_name/123"
		result, err := validator.MakeMetricValuesFromURL(url)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "counter", result.MetricType)
		assert.Equal(t, "metric_name", result.MetricName)
		assert.Equal(t, "123", result.MetricValue)
	})

	t.Run("Неверный URL-адрес - недостаточно частей", func(t *testing.T) {
		url := "/update/counter/metric_name"
		result, err := validator.MakeMetricValuesFromURL(url)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotCorrectURL, err)
	})

	// Add more test cases as needed
}

func TestMetricValidator_GetValidType(t *testing.T) {
	validator := New()

	t.Run("Допустимый тип - Counter", func(t *testing.T) {
		result := validator.GetValidType("counter")
		assert.Equal(t, entity.Counter, result)
	})

	t.Run("Допустимый тип - Gauge", func(t *testing.T) {
		result := validator.GetValidType("gauge")
		assert.Equal(t, entity.Gauge, result)
	})

	t.Run("Недопустимый тип", func(t *testing.T) {
		result := validator.GetValidType("invalid_type")
		assert.Equal(t, entity.NoType, result)
	})
}

func TestMetricValidator_GetTypeGaugeValue(t *testing.T) {
	validator := New()

	t.Run("Допустимое значение Gauge", func(t *testing.T) {
		result, err := validator.GetTypeGaugeValue("123.45")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 123.45, *result)
	})

	t.Run("Недопустимое значение Gauge", func(t *testing.T) {
		result, err := validator.GetTypeGaugeValue("invalid_value")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotCorrectValue, err)
	})
}

func TestMetricValidator_GetTypeCounterValue(t *testing.T) {
	validator := New()

	t.Run("Допустимое значение Home", func(t *testing.T) {
		result, err := validator.GetTypeCounterValue("42")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(42), *result)
	})

	t.Run("Недопустимое значение Home", func(t *testing.T) {
		result, err := validator.GetTypeCounterValue("invalid_value")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotCorrectValue, err)
	})
}
