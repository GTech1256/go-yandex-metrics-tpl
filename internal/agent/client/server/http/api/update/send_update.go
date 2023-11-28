package update

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api"

	netHTTP "net/http"
)

type RequestError struct {
}

func (s update) SendUpdate(ctx context.Context, updateDto dto.Update) error {
	requestURL := getRequestURL(s.BaseURL, &updateDto)

	req, err := s.HTTPClient.NewRequest(netHTTP.MethodPost, requestURL, nil)
	if err != nil {
		s.logger.Infof("client: could not create request: %s\n", err)
		return err
	}

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		s.logger.Errorf("client: error making service request: %w\n", err)

		return api.ErrRequestDo
	}

	defer res.Body.Close()

	s.logger.Infof("%d %v \n", res.StatusCode, requestURL)

	if res.StatusCode != netHTTP.StatusOK {
		s.logger.Errorf("%v v", api.ErrInvalidResponseStatus, res.StatusCode)

		return api.ErrInvalidResponseStatus
	}

	return nil
}

func getRequestURL(baseURL string, updateDto *dto.Update) string {
	return fmt.Sprintf("%v/update/%v/%v/%v", baseURL, updateDto.Type, updateDto.Name, updateDto.Value)
}
