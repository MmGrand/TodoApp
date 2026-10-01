package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/MmGrand/TodoApp/internal/core/transport/http/middleware"
)

type ApiVersion string

const (
	ApiVersion1 ApiVersion = "v1"
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	Middleware []core_http_middleware.Middleware
}

func NewApiVersionRouter(
	apiVersion ApiVersion,
	Middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		Middleware: Middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.Middleware...,
	)
}
