package guard

import (
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	metricvalidator "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service/metric_validator"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
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
		getMockMetricValidator func() *metricvalidator.MockMetricValidator
		getMockLogger          func() *logging.LoggerMock
	}{
		{
			name:    "Gauge/Title/10.10 Success",
			urlPath: "/update/gauge/title/10.10",
			method:  http.MethodPost,
			getMockMetricValidator: func() *metricvalidator.MockMetricValidator {
				// TODO: Поправить путь на "/update/counter/title/10"
				expURL := ""
				mockMetricValidator := new(metricvalidator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expURL).Return(&entity2.MetricFields{
					MetricType:  string(entity2.Gauge),
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
			getMockMetricValidator: func() *metricvalidator.MockMetricValidator {
				// TODO: Поправить путь на "/update/counter/title/10"
				expURL := ""
				mockMetricValidator := new(metricvalidator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expURL).Return(&entity2.MetricFields{
					MetricType:  string(entity2.Counter),
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
			getMockMetricValidator: func() *metricvalidator.MockMetricValidator {
				// TODO: Поправить путь
				expURL := ""
				mockMetricValidator := new(metricvalidator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expURL).Return(&entity2.MetricFields{}, metricvalidator.ErrNotCorrectURL)

				return mockMetricValidator
			},
			getMockLogger: func() *logging.LoggerMock {
				mockLogger := new(logging.LoggerMock)
				mockLogger.On("Error", metricvalidator.ErrNotCorrectURL)

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
			getMockMetricValidator: func() *metricvalidator.MockMetricValidator {
				// TODO: Поправить путь
				expURL := ""
				mockMetricValidator := new(metricvalidator.MockMetricValidator)
				mockMetricValidator.On("MakeMetricValuesFromURL", expURL).Return(&entity2.MetricFields{}, metricvalidator.ErrNotCorrectName)

				return mockMetricValidator
			},
			getMockLogger: func() *logging.LoggerMock {
				mockLogger := new(logging.LoggerMock)
				mockLogger.On("Error", metricvalidator.ErrNotCorrectName)

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
			getMockMetricValidator: func() *metricvalidator.MockMetricValidator {
				mockMetricValidator := new(metricvalidator.MockMetricValidator)

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
