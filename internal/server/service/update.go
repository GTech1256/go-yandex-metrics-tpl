package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http/update/interface"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var (
	// ErrNotCorrectURL возвращает ошибку, если формат url не подошел
	ErrNotCorrectURL   = errors.New("not correct url")
	ErrNotCorrectType  = errors.New("not correct metricType")
	ErrNotCorrectValue = errors.New("not correct metricValue")
)

type updateService struct {
	logger  *logrus.Entry
	storage metric2.Storage
}

func NewUpdateService(logger *logrus.Entry, storage metric2.Storage) updateInterface.Service {
	return &updateService{logger: logger, storage: storage}
}

func (u updateService) GetMetric(ctx context.Context, url string) (*entity.Metric, error) {
	const count = 5
	splitted := strings.Split(url, "/")
	u.logger.Info(splitted)

	if len(splitted) != count {
		u.logger.Info()
		u.logger.Info("len(splitted) != count", len(splitted), splitted)
		fmt.Printf("%+v", splitted)

		return nil, ErrNotCorrectURL
	}

	metricTypUnsafe, metricName, metricValueUnsafe := splitted[2], splitted[3], splitted[4]

	metricType := getValidType(metricTypUnsafe)
	if metricType == entity.NoType {
		return nil, ErrNotCorrectType
	}

	var metricValue interface{}
	var err error
	if metricType == entity.Gauge {
		metricValue, err = getTypeGaugeValue(metricValueUnsafe)
		if err != nil {
			return nil, ErrNotCorrectValue
		}

	} else if metricType == entity.Counter {
		metricValue, err = getTypeCounterValue(metricValueUnsafe)
		if err != nil {
			return nil, ErrNotCorrectValue
		}
	} else {
		return nil, ErrNotCorrectValue
	}

	return &entity.Metric{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	}, nil
}

func getValidType(metricType string) entity.Type {
	switch entity.Type(metricType) {
	case entity.Gauge:
		return entity.Gauge
	case entity.Counter:
		return entity.Counter
	default:
		return entity.NoType
	}
}

func getTypeGaugeValue(metricValueUnsafe string) (*float64, error) {
	metricValue, err := strconv.ParseFloat(metricValueUnsafe, 64)
	if err != nil {
		return nil, ErrNotCorrectValue
	}

	return &metricValue, nil
}

func getTypeCounterValue(metricValueUnsafe string) (*int64, error) {
	metricValue, err := strconv.ParseInt(metricValueUnsafe, 10, 64)
	if err != nil {
		return nil, ErrNotCorrectValue
	}

	return &metricValue, nil
}

func (u updateService) SaveGaugeMetric(ctx context.Context, metric *entity.MetricGauge) error {
	return u.storage.SaveGauge(ctx, metric)
}
func (u updateService) SaveCounterMetric(ctx context.Context, metric *entity.MetricCounter) error {
	return u.storage.SaveCounter(ctx, metric)
}
