package server

import (
	"io"
	"net/http"
)

type HttpClient interface {
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
	//defaultClient
}

type httpClient struct {
	//NewRequest func(method, url string, body io.Reader) (*http.Request, error)
	//Do func(req *http.Request) (*http.Response, error)
}

func (h httpClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func (h httpClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
