package updates

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/http/api/update/converter"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/gzip"

	netHTTP "net/http"
)

func (s updates) SendUpdates(ctx context.Context, updateDto []*dto.Update) error {
	requestURL := getRequestURLJSON(s.BaseURL)

	ms := make([]*converter.Metrics, 0)

	for _, el := range updateDto {
		m, err := converter.UpdateDTOToMetrics(el)
		if err != nil {
			return err
		}

		ms = append(ms, m)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(ms)
	if err != nil {
		return err
	}

	bufCompress, err := gzip.Compress(buf.Bytes())
	if err != nil {
		return err
	}

	req, err := s.HTTPClient.NewRequest(netHTTP.MethodPost, requestURL, bytes.NewReader(bufCompress))
	if err != nil {
		s.logger.Infof("client: could not create request: %s\n", err)
		return err
	}

	req.Header.Add("Content-Encoding", "gzip")

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		s.logger.Infof("client: error making service request: %s\n", err)
		return err
	}
	defer res.Body.Close()

	s.logger.Infof("%d %v %+v \n", res.StatusCode, requestURL, updateDto)

	return nil
}

func getRequestURLJSON(baseURL string) string {
	return fmt.Sprintf("%v/updates", baseURL)
}
