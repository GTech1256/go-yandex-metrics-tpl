package metricvalidator

import (
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"strconv"
	"strings"
)

const updatePartLength = 8 // len("/update/")

type metricValidator struct {
}

func New() *metricValidator {
	return &metricValidator{}
}

// MakeMetricValuesFromURL На вход ожидается /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (m metricValidator) MakeMetricValuesFromURL(url string) (*entity2.MetricFields, error) {
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
	if validType == entity2.NoType {
		return nil, ErrNotCorrectType
	}

	return &entity2.MetricFields{
		MetricType:  metricType,
		MetricName:  metricName,
		MetricValue: metricValue,
	}, nil

}

func (m metricValidator) GetValidType(metricType string) entity2.Type {
	switch entity2.Type(metricType) {
	case entity2.Gauge:
		return entity2.Gauge
	case entity2.Counter:
		return entity2.Counter
	default:
		return entity2.NoType
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
