package service

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/adapters/http/update/interface"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
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
