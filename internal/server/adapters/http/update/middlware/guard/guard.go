package guard

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/util"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	ExpectContentType = "text/plain"
	ExpectMethod      = http.MethodPost
)

func WithMetricGuarding(next http.Handler, logger *logrus.Entry) http.Handler {
	guardFn := func(rw http.ResponseWriter, req *http.Request) {
		isCorrectMethod := req.Method == ExpectMethod
		if !isCorrectMethod {
			logger.Info("Allow Only Method POST, Got: ", req.Method)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := util.MakeMetricValuesFromUrl(req.RequestURI)

		if err == service.ErrNotCorrectName {
			logger.Info(service.ErrNotCorrectName)
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			logger.Info(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		//if err == service.ErrNotCorrectType {
		//	rw.WriteHeader(http.StatusBadRequest)
		//	return
		//}

		next.ServeHTTP(rw, req)
	}
	return http.HandlerFunc(guardFn)
}
