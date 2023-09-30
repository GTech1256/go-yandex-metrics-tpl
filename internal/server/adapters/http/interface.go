package http

import "net/http"

type Handler interface {
	Register(router *http.ServeMux)
}
