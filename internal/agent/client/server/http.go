package server

import (
	"io"
	"net/http"
)

type HTTPClient interface {
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
}

func (h httpClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func (h httpClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
