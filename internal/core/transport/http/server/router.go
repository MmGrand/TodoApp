package core_http_server

import (
	"fmt"
	"net/http"

	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_http_response "github.com/MmGrand/TodoApp/internal/core/transport/http/response"
)

type APIVersion string

const (
	APIVersion1 APIVersion = "v1"
)

type APIVersionRouter struct {
	mux        *http.ServeMux
	apiVersion APIVersion
}

func NewAPIVersionRouter(apiVersion APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		mux:        http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.mux.Handle(pattern, route.Handler)
	}
}

func (r *APIVersionRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, pattern := r.mux.Handler(req)
	if pattern != "" {
		r.mux.ServeHTTP(w, req)

		return
	}

	interceptor := &routeErrorInterceptor{ResponseWriter: w}
	handler.ServeHTTP(interceptor, req)

	if interceptor.statusCode == 0 {
		return
	}

	msg := "route not found"
	if interceptor.statusCode == http.StatusMethodNotAllowed {
		msg = "method not allowed"
	}

	log := core_logger.FromContext(req.Context())
	core_http_response.NewHTTPResponseHandler(log, w).JSONResponse(
		core_http_response.ErrorResponse{
			Error:   fmt.Sprintf("%s %s: %s", req.Method, req.URL.Path, msg),
			Message: msg,
		},
		interceptor.statusCode,
	)
}

type routeErrorInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (i *routeErrorInterceptor) WriteHeader(statusCode int) {
	if statusCode == http.StatusNotFound || statusCode == http.StatusMethodNotAllowed {
		i.statusCode = statusCode

		return
	}

	i.ResponseWriter.WriteHeader(statusCode)
}

func (i *routeErrorInterceptor) Write(b []byte) (int, error) {
	if i.statusCode != 0 {
		return len(b), nil
	}

	return i.ResponseWriter.Write(b)
}
