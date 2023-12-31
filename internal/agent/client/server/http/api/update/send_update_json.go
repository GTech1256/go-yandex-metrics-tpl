package update

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api/update/converter"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/gzip"

	netHTTP "net/http"
)

var (
	ErrCompression = errors.New("при сжатии данных произошла ошибка")
	ErrMarshal     = errors.New("при переводе структуры в json произошла ошибка")
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
		s.logger.Error(err)
		return ErrMarshal
	}

	bufCompress, err := gzip.Compress(buf.Bytes())
	if err != nil {
		s.logger.Error(err)
		return ErrCompression
	}

	req, err := s.HTTPClient.NewRequest(netHTTP.MethodPost, requestURL, bytes.NewReader(bufCompress))
	if err != nil {
		s.logger.Infof("client: could not create request: %s\n", err)
		return api.ErrRequestInitiate
	}

	req.Header.Add("Content-Encoding", "gzip")

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		s.logger.Error("client: error making service request: %w\n", err)

		return api.ErrRequestDo
	}
	defer res.Body.Close()

	s.logger.Infof("%d %v %+v \n", res.StatusCode, requestURL, updateDto)

	if res.StatusCode != netHTTP.StatusOK {
		s.logger.Errorf("%v v", api.ErrInvalidResponseStatus, res.StatusCode)

		return api.ErrInvalidResponseStatus
	}

	return nil
}

func getRequestURLJSON(baseURL string) string {
	return fmt.Sprintf("%v/update/", baseURL)
}
