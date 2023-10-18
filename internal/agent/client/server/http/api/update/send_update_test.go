package update

import (
	"context"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	serverHttp "github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http"
	logging2 "github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func Test_client_Post(t *testing.T) {
	host := "http://example.com"
	mockLogger := new(logging2.LoggerMock)
	mockHTTPClient := new(serverHttp.MockHTTPClient)
	client2 := New(mockHTTPClient, host, mockLogger)
	// Подготовка данных для теста
	update := dto.Update{
		Type:  "type",
		Name:  "name",
		Value: "value",
	}

	requestURL := getRequestURL(host, &update)
	
	// Устанавливаем ожидаемые вызовы для мока ClientHTTP
	mockHTTPClient.On("NewRequest", http.MethodPost, requestURL, mock.Anything).Return(&http.Request{}, nil)
	bodyMock := new(serverHttp.BodyMock)
	mockHTTPClient.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: bodyMock}, nil)

	mockLogger.On("Infof", mock.Anything, mock.Anything, mock.Anything)

	// Выполняем функцию
	err := client2.SendUpdate(context.Background(), update)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем ожидаемые вызовы
	mockHTTPClient.AssertExpectations(t)
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
