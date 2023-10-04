package util

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service"
	"strconv"
	"strings"
)

const updatePartLength = 8 // len("/update/")

// MakeMetricValuesFromURL На вход ожидается /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func MakeMetricValuesFromURL(url string) (*entity.MetricFields, error) {
	if len(url) < 8 {
		return nil, service.ErrNotCorrectURL
	}

	urlWithoutUpdate := url[updatePartLength:]
	splitted := strings.Split(urlWithoutUpdate, "/")

	switch len(splitted) {
	case 0: // ничего нет
	case 2: // <ТИП_МЕТРИКИ> <ИМЯ_МЕТРИКИ>
		isNameEmpty := len(splitted[1]) == 0
		if isNameEmpty {
			return nil, service.ErrNotCorrectName
		}

		return nil, service.ErrNotCorrectURL
	case 1: // <ТИП_МЕТРИКИ> : Нет имени
		return nil, service.ErrNotCorrectName

	}

	metricType, metricName, metricValue := splitted[0], splitted[1], splitted[2]

	if len(metricName) == 0 {
		return nil, service.ErrNotCorrectName
	}

	validType := GetValidType(metricType)
	if validType == entity.NoType {
		return nil, service.ErrNotCorrectType
	}

	return &entity.MetricFields{
		MetricType:  metricType,
		MetricName:  metricName,
		MetricValue: metricValue,
	}, nil

}

func GetValidType(metricType string) entity.Type {
	switch entity.Type(metricType) {
	case entity.Gauge:
		return entity.Gauge
	case entity.Counter:
		return entity.Counter
	default:
		return entity.NoType
	}
}

func GetTypeGaugeValue(metricValueUnsafe string) (*float64, error) {
	metricValue, err := strconv.ParseFloat(metricValueUnsafe, 64)
	if err != nil {
		return nil, service.ErrNotCorrectValue
	}

	return &metricValue, nil
}

func GetTypeCounterValue(metricValueUnsafe string) (*int64, error) {
	metricValue, err := strconv.ParseInt(metricValueUnsafe, 10, 64)
	if err != nil {
		return nil, service.ErrNotCorrectValue
	}

	return &metricValue, nil
}
