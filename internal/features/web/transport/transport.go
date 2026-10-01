package web_transport_http

import (
	core_http_server "github.com/MmGrand/TodoApp/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHTTPHandler(
	WwebService WebService,
) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: WwebService,
	}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
