package logging

import (
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// our http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter // compose original http.ResponseWriter
		responseData        *responseData
	}
)

func WithLogging(h http.Handler, logger logging2.Logger) http.Handler {
	loggingFn := func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: rw, // compose original service.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lrw, req) // inject our implementation of service.ResponseWriter

		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"uri":      req.RequestURI,
			"method":   req.Method,
			"status":   responseData.status,
			"duration": duration,
			"size":     responseData.size,
		}).Info("request completed")
	}
	return http.HandlerFunc(loggingFn)
}
