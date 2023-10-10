package value

import (
	"fmt"
	updateinterface "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/http/update/interface"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/service"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/lib/ptr"
	logging "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValueHandler(t *testing.T) {
	t.Run("Успешный GetMetricValue", func(t *testing.T) {
		mockLogger := new(logging.LoggerMock)
		mockService := new(service.MockService)
		mockLogger.On("Error").Return(nil)

		// Тестовый маршрутизатор
		router := chi.NewRouter()
		h := NewHandler(mockLogger, mockService, nil) // Передаем мок сервиса
		h.Register(router)

		req := httptest.NewRequest("GET", "/value/testType/testName", nil)
		rec := httptest.NewRecorder()

		mockService.On("GetMetricValue", mock.Anything, mock.MatchedBy(func(dto *updateinterface.GetMetricValueDto) bool {
			return dto.Type == "testType" && dto.Name == "testName"
		})).Return(ptr.StrPtr("123.45"), nil)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "123.45", rec.Body.String())
		mockLogger.AssertNotCalled(t, "Error")
		mockService.AssertExpectations(t)
	})

	t.Run("Неуспешный GetMetricValue", func(t *testing.T) {
		mockLogger := new(logging.LoggerMock)
		mockService := new(service.MockService)
		// Тестовый маршрутизатор
		router := chi.NewRouter()
		h := NewHandler(mockLogger, mockService, nil) // Передаем мок сервиса
		h.Register(router)
		mockLogger.On("Error", fmt.Errorf("Ошибка")).Return(nil)

		req := httptest.NewRequest("GET", "/value/testType/testName", nil)
		rec := httptest.NewRecorder()

		mockService.On("GetMetricValue", mock.Anything, mock.MatchedBy(func(dto *updateinterface.GetMetricValueDto) bool {
			return dto.Type == "testType" && dto.Name == "testName"
		})).Return(nil, fmt.Errorf("Ошибка"))

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Empty(t, rec.Body.String())
		//mockLogger.AssertCalled(t, "Error")
		mockLogger.AssertExpectations(t)
		mockService.AssertExpectations(t)
	})
}
