package http

import (
	"io"
	netHTTP "net/http"
)

type httpClient struct {
}

func New() ClientHTTP {
	return &httpClient{}
}

func (h httpClient) NewRequest(method, url string, body io.Reader) (*netHTTP.Request, error) {
	return netHTTP.NewRequest(method, url, body)
}

func (h httpClient) Do(req *netHTTP.Request) (*netHTTP.Response, error) {
	return netHTTP.DefaultClient.Do(req)
}
