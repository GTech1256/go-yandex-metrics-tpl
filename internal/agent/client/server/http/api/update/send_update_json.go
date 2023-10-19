package update

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api/update/converter"

	netHTTP "net/http"
)

func (s update) SendUpdateJSON(ctx context.Context, updateDto dto.Update) error {
	requestURL := getRequestURLJSON(s.BaseURL)

	m, err := converter.UpdateDTOToMetrics(&updateDto)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err = encoder.Encode(m)
	if err != nil {
		return err
	}

	req, err := s.HTTPClient.NewRequest(netHTTP.MethodPost, requestURL, &buf)
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

func getRequestURLJSON(baseURL string) string {
	return fmt.Sprintf("%v/update/", baseURL)
}
