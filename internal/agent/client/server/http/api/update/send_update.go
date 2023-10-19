package update

import (
	"context"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"

	netHTTP "net/http"
)

func (s update) SendUpdate(ctx context.Context, updateDto dto.Update) error {
	requestURL := getRequestURL(s.BaseURL, &updateDto)

	fmt.Println("REST REST REST")
	fmt.Println("REST REST REST")
	fmt.Println("REST REST REST")
	req, err := s.HTTPClient.NewRequest(netHTTP.MethodPost, requestURL, nil)
	if err != nil {
		s.logger.Infof("client: could not create request: %s\n", err)
		return err
	}

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		s.logger.Infof("client: error making service request: %s\n", err)
		return err
	}
	defer res.Body.Close()

	s.logger.Infof("%d %v \n", res.StatusCode, requestURL)

	return nil
}

func getRequestURL(baseURL string, updateDto *dto.Update) string {
	return fmt.Sprintf("%v/update/%v/%v/%v", baseURL, updateDto.Type, updateDto.Name, updateDto.Value)
}
