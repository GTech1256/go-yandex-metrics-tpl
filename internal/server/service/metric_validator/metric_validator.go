package metricvalidator

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"strconv"
	"strings"
)

const updatePartLength = 8 // len("/update/")

type MetricValidator interface {
	MakeMetricValuesFromURL(url string) (*entity.MetricFields, error)
	GetValidType(metricType string) entity.Type
	GetTypeGaugeValue(metricValueUnsafe string) (*float64, error)
	GetTypeCounterValue(metricValueUnsafe string) (*int64, error)
}

type metricValidator struct {
}

func New() MetricValidator {
	return &metricValidator{}
}

// MakeMetricValuesFromURL На вход ожидается /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (m metricValidator) MakeMetricValuesFromURL(url string) (*entity.MetricFields, error) {
	if len(url) < 8 {
		return nil, ErrNotCorrectURL
	}

	urlWithoutUpdate := url[updatePartLength:]
	splitted := strings.Split(urlWithoutUpdate, "/")

	switch len(splitted) {
	case 0: // ничего нет
	case 2: // <ТИП_МЕТРИКИ> <ИМЯ_МЕТРИКИ>
		isNameEmpty := len(splitted[1]) == 0
		if isNameEmpty {
			return nil, ErrNotCorrectName
		}

		return nil, ErrNotCorrectURL
	case 1: // <ТИП_МЕТРИКИ> : Нет имени
		return nil, ErrNotCorrectName

	}

	metricType, metricName, metricValue := splitted[0], splitted[1], splitted[2]

	if len(metricName) == 0 {
		return nil, ErrNotCorrectName
	}

	validType := m.GetValidType(metricType)
	if validType == entity.NoType {
		return nil, ErrNotCorrectType
	}

	return &entity.MetricFields{
		MetricType:  metricType,
		MetricName:  metricName,
		MetricValue: metricValue,
	}, nil

}

func (m metricValidator) GetValidType(metricType string) entity.Type {
	switch entity.Type(metricType) {
	case entity.Gauge:
		return entity.Gauge
	case entity.Counter:
		return entity.Counter
	default:
		return entity.NoType
	}
}

func (m metricValidator) GetTypeGaugeValue(metricValueUnsafe string) (*float64, error) {
	metricValue, err := strconv.ParseFloat(metricValueUnsafe, 64)
	if err != nil {
		return nil, ErrNotCorrectValue
	}

	return &metricValue, nil
}

func (m metricValidator) GetTypeCounterValue(metricValueUnsafe string) (*int64, error) {
	metricValue, err := strconv.ParseInt(metricValueUnsafe, 10, 64)
	if err != nil {
		return nil, ErrNotCorrectValue
	}

	return &metricValue, nil
}
