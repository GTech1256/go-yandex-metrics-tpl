package server

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func Test_client_Post(t *testing.T) {
	host := "http://example.com"

	mockLogger := new(logging2.LoggerMock)
	mockLogger.On("Infof", mock.Anything, mock.Anything, mock.Anything)

	client2 := New(host, mockLogger)

	// Создаем мок для HTTPClient
	httpClientMock := &MockHTTPClient{}
	client2.(*client).httpClient = httpClientMock

	// Подготовка данных для теста
	update := dto.Update{
		Type:  "type",
		Name:  "name",
		Value: "value",
	}

	requestURL := getRequestURL(host, &update)

	bodyMock := new(BodyMock)

	// Устанавливаем ожидаемые вызовы для мока HTTPClient
	httpClientMock.On("NewRequest", http.MethodPost, requestURL, mock.Anything).Return(&http.Request{}, nil)
	httpClientMock.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: bodyMock}, nil)

	// Выполняем функцию
	err := client2.Post(context.Background(), update)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем ожидаемые вызовы
	httpClientMock.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func Test_getRequestURL(t *testing.T) {
	type args struct {
		host      string
		updateDto *dto.Update
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				host: "http://example.com",
				updateDto: &dto.Update{
					Type:  "testType",
					Name:  "testName",
					Value: "testValue",
				},
			},
			want: "http://example.com/update/testType/testName/testValue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getRequestURL(tt.args.host, tt.args.updateDto), "getRequestURL(%v, %v)", tt.args.host, tt.args.updateDto)
		})
	}
}
