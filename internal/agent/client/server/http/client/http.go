package client

import (
	"bytes"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"io"
	netHTTP "net/http"
	"time"
)

type httpClient struct {
	HashKey *string
	logger  logging.Logger
}

func New(HashKey *string, logger logging.Logger) ClientHTTP {
	return &httpClient{
		HashKey: HashKey,
		logger:  logger,
	}
}

func (h httpClient) NewRequest(method, url string, body io.Reader) (*netHTTP.Request, error) {
	time.Sleep(time.Second)
	r, err := netHTTP.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if h.HashKey != nil && body != nil {
		bodyByte, err := io.ReadAll(body)

		r.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
		if err != nil {
			return nil, err
		}

		hash, err := h.getEncryptBodyHash(bodyByte)
		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		hashEncode := h.getEncode(hash)

		r.Header.Set("HashSHA256", hashEncode)
	}

	return r, err
}

func (h httpClient) Do(req *netHTTP.Request) (*netHTTP.Response, error) {
	return netHTTP.DefaultClient.Do(req)
}
