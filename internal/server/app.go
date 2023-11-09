package server

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/composition"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/config"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	home "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/gzip"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/middlware/logging"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update"
	updateInterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/rest/counter"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/update/rest/gauge"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/updates"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/rest/value"
	file2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/repository/metric/file"
	metric2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_loader"
	metricValidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type App struct {
	logger          logging2.Logger
	router          *chi.Mux
	cfg             *config.Config
	metricService   MetricService
	metricValidator MetricValidator
}

type MetricLoaderService interface {
	StartMetricsToDiskInterval(ctx context.Context, interval time.Duration)
	LoadMetricsFromDisk(ctx context.Context) ([]*file2.MetricJSON, error)
	SaveMetricToDisk(ctx context.Context, mj *file2.MetricJSON) error
}

type MetricService interface {
	SaveGaugeMetric(ctx context.Context, metric *entity2.MetricFields) error
	SaveCounterMetric(ctx context.Context, metric *entity2.MetricFields) error
	GetMetricValue(ctx context.Context, metric *updateInterface.GetMetricValueDto) (*string, error)
	GetMetrics(ctx context.Context) (*metric.AllMetrics, error)
	SaveMetricJSONs(ctx context.Context, metrics []*entity2.MetricJSON) error
	SaveMetricJSON(ctx context.Context, metric *entity2.MetricJSON) error
}

type MetricValidator interface {
	MakeMetricValuesFromURL(url string) (*entity2.MetricFields, error)
	GetValidType(metricType string) entity2.Type
	GetTypeGaugeValue(metricValueUnsafe string) (*float64, error)
	GetTypeCounterValue(metricValueUnsafe string) (*int64, error)
}

func New(cfg *config.Config, logger logging2.Logger) (*App, error) {
	router := chi.NewRouter()
	router.Use(gzip.WithGzip)

	storage, err := composition.NewStorageComposition(cfg, logger, router)
	if err != nil {
		logger.Error(err)
	}

	var metricLoaderService MetricLoaderService = nil
	// пустое значение отключает функцию записи на диск
	if cfg.GetIsEnabledFileWrite() && !cfg.GetIsEnabledSQLStore() {
		fileStorage, err := file2.NewFileStorage(*cfg.FileStoragePath)
		if err != nil {
			logger.Error("Ошибка инцилизации fileStorage ", err)
			panic(err)
		}

		metricLoaderService = metricloader.NewMetricLoaderService(logger, fileStorage, storage)
	}

	validator := metricValidator.New()
	metricService := metric2.NewMetricService(logger, storage, validator, metricLoaderService, cfg)

	app := &App{
		logger:          logger,
		router:          router,
		metricService:   metricService,
		cfg:             cfg,
		metricValidator: validator,
	}

	app.handlersRegister()

	return app, nil
}

func (a App) handlersRegister() {
	a.logger.Info("Регистрация /update/counter/ Router")
	updateCounterHandler := counter.NewHandler(a.logger, a.metricService, a.metricValidator)
	updateCounterHandler.Register(a.router)

	a.logger.Info("Регистрация /update/gauge/ Router")
	updateGaugeHandler := gauge.NewHandler(a.logger, a.metricService, a.metricValidator)
	updateGaugeHandler.Register(a.router)

	a.logger.Info("Регистрация /update/* Router")
	updateHandler := update.NewHandler(a.logger, a.metricService, a.metricValidator)
	updateHandler.Register(a.router)

	a.logger.Info("Регистрация /updates Router")
	updatesHandler := updates.NewHandler(a.logger, a.metricService, a.metricValidator)
	updatesHandler.Register(a.router)

	a.logger.Info("Регистрация /value/ Router")
	valueHandler := value.NewHandler(a.logger, a.metricService, a.metricValidator)
	valueHandler.Register(a.router)

	a.logger.Info("Регистрация / Router")
	homeHandler := home.NewHandler(a.logger, a.metricService)
	homeHandler.Register(a.router)

	a.logger.Info(fmt.Sprintf("Start Listen Port %v", *a.cfg.Port))
	log.Fatal(http.ListenAndServe(*a.cfg.Port, logging.WithLogging(a.router, a.logger)))
}
