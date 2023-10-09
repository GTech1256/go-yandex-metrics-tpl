package server

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"net/http"
)

func (s client) Post(ctx context.Context, updateDto dto.Update) error {
	requestURL := getRequestURL(s.host, &updateDto)

	req, err := s.httpClient.NewRequest(http.MethodPost, requestURL, nil)

	if err != nil {
		s.logger.Infof("client: could not create request: %s\n", err)
		return err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Infof("client: error making service request: %s\n", err)
		return err
	}
	defer res.Body.Close()

	s.logger.Infof("%d %v \n", res.StatusCode, requestURL)

	return nil
}

func getRequestURL(host string, updateDto *dto.Update) string {
	return fmt.Sprintf("%v/update/%v/%v/%v", host, updateDto.Type, updateDto.Name, updateDto.Value)
}
