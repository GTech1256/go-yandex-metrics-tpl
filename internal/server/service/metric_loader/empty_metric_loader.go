package metricloader

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/file"
	"time"
)

type emptyMetricLoaderService struct {
}

func (e emptyMetricLoaderService) StartMetricsToDiskInterval(ctx context.Context, interval time.Duration) {
}

func (e emptyMetricLoaderService) LoadMetricsFromDisk(ctx context.Context) ([]*file.MetricJSON, error) {
	return nil, nil
}

func (e emptyMetricLoaderService) SaveMetricToDisk(ctx context.Context, mj *file.MetricJSON) error {
	return nil
}

func NewEmptyMetricLoaderService() *emptyMetricLoaderService {
	return &emptyMetricLoaderService{}
}
