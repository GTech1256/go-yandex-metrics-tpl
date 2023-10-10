package guard

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"net/http"
)

const (
	ExpectMethod = http.MethodPost
)

func WithMetricGuarding(next http.Handler, logger logging2.Logger, validator metricValidator.MetricValidator) http.HandlerFunc {
	guardFn := func(rw http.ResponseWriter, req *http.Request) {
		isCorrectMethod := req.Method == ExpectMethod
		if !isCorrectMethod {
			logger.Error("Allow Only Method POST, Got: ", req.Method)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := validator.MakeMetricValuesFromURL(req.RequestURI)

		if err == metricValidator.ErrNotCorrectName {
			logger.Error(metricValidator.ErrNotCorrectName)
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			logger.Error(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(rw, req)
	}
	return guardFn
}
