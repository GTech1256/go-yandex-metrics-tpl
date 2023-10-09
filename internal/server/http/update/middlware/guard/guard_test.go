package guard

import (
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithMetricGuarding(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name                   string
		want                   want
		urlPath                string
		method                 string
		getMockMetricValidator func() *metric_validator.MockMetricValidator
		getMockLogger          func() *logging.LoggerMock
	}{
		{
			name:    "Gauge/Title/10.10 Success",
			urlPath: "/update/gauge/title/10.10",
			method:  http.MethodPost,
			getMockMetricValidator: func() *metric_validator.MockMetricValidator {
				// TODO: Поправить путь на "/update/counter/title/10"
				expUrl := ""
				mockMetricValidator := new(metric_validator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expUrl).Return(&entity.MetricFields{
					MetricType:  string(entity.Gauge),
					MetricName:  "title",
					MetricValue: "10.10",
				}, nil)

				return mockMetricValidator
			},
			getMockLogger: func() *logging.LoggerMock {
				mockLogger := new(logging.LoggerMock)

				return mockLogger
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/text",
			},
		},
		{
			name:    "Counter/Title/10 Success",
			urlPath: "/update/counter/title/10",
			method:  http.MethodPost,
			getMockMetricValidator: func() *metric_validator.MockMetricValidator {
				// TODO: Поправить путь на "/update/counter/title/10"
				expUrl := ""
				mockMetricValidator := new(metric_validator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expUrl).Return(&entity.MetricFields{
					MetricType:  string(entity.Counter),
					MetricName:  "title",
					MetricValue: "10",
				}, nil)

				return mockMetricValidator
			},
			getMockLogger: func() *logging.LoggerMock {
				mockLogger := new(logging.LoggerMock)

				return mockLogger
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/text",
			},
		},
		{
			name:    "Fail Not valid path",
			urlPath: "/update/gauge",
			method:  http.MethodPost,
			getMockMetricValidator: func() *metric_validator.MockMetricValidator {
				// TODO: Поправить путь
				expUrl := ""
				mockMetricValidator := new(metric_validator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expUrl).Return(&entity.MetricFields{}, metric_validator.ErrNotCorrectURL)

				return mockMetricValidator
			},
			getMockLogger: func() *logging.LoggerMock {
				mockLogger := new(logging.LoggerMock)
				mockLogger.On("Error", metric_validator.ErrNotCorrectURL)

				return mockLogger
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/text",
			},
		},
		{
			name:    "Fail Not valid name",
			urlPath: "/update/",
			method:  http.MethodPost,
			getMockMetricValidator: func() *metric_validator.MockMetricValidator {
				// TODO: Поправить путь
				expUrl := ""
				mockMetricValidator := new(metric_validator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expUrl).Return(&entity.MetricFields{}, metric_validator.ErrNotCorrectName)

				return mockMetricValidator
			},
			getMockLogger: func() *logging.LoggerMock {
				mockLogger := new(logging.LoggerMock)
				mockLogger.On("Error", metric_validator.ErrNotCorrectName)

				return mockLogger
			},
			want: want{
				code:        http.StatusNotFound,
				contentType: "application/text",
			},
		},
		{
			name:    "Fail Not HTTP.Method",
			urlPath: "/update/gauge/title/10.10",
			method:  http.MethodGet,
			getMockMetricValidator: func() *metric_validator.MockMetricValidator {
				mockMetricValidator := new(metric_validator.MockMetricValidator)

				return mockMetricValidator
			},
			getMockLogger: func() *logging.LoggerMock {
				mockLogger := new(logging.LoggerMock)
				mockLogger.On("Error", "Allow Only Method POST, Got: ", http.MethodGet)

				return mockLogger
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/text",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			mockLogger := tt.getMockLogger()
			mockMetricValidator := tt.getMockMetricValidator()
			handler := WithMetricGuarding(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusOK)
			}), mockLogger, mockMetricValidator)

			req, err := http.NewRequest(tt.method, tt.urlPath, nil)
			assert.NoError(t, err)

			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tt.want.code, recorder.Code)
		})
	}
}
