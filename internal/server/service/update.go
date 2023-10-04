package service

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/sirupsen/logrus"
)

type updateService struct {
	logger  *logrus.Entry
	storage metric2.Storage
}

func NewUpdateService(logger *logrus.Entry, storage metric2.Storage) updateInterface.Service {
	return &updateService{logger: logger, storage: storage}
}

func (u updateService) SaveGaugeMetric(ctx context.Context, metric *entity.MetricGauge) error {
	return u.storage.SaveGauge(ctx, metric)
}
func (u updateService) SaveCounterMetric(ctx context.Context, metric *entity.MetricCounter) error {
	return u.storage.SaveCounter(ctx, metric)
}
