package logging

import (
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	body        []byte
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b

	return rw.ResponseWriter.Write(b)
}

func WithLogging(h http.Handler, logger logging2.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(
					"err", err,
					"trace", debug.Stack(),
				)
			}
		}()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		h.ServeHTTP(wrapped, r)

		//var buf bytes.Buffer
		//
		//_, err := buf.ReadFrom(r.Body)
		//if err != nil {
		//	logger.Error(err)
		//	return
		//}
		//
		//r.Body.Close()

		//fmt.Println("JIHJO", buf.String(), buf.Bytes(), r.Body)
		//fmt.Println("JIHJO", json.NewEncoder(NewReade(wrapped.body)), string(wrapped.body))
		//body, err := io.ReadAll(r.Body)
		//if err != nil {
		//	logger.Error(err)
		//}

		logger.WithFields(logrus.Fields{
			"status":   wrapped.status,
			"path":     r.URL.EscapedPath(),
			"method":   r.Method,
			"duration": time.Since(start),
			"length":   r.ContentLength,
			//"inputBody": buf.String(),
			//"outBody":   string(wrapped.body),
		}).Info("request completed")

	}

	return http.HandlerFunc(fn)
}
