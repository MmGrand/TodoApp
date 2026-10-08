package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	"go.uber.org/zap"
)

const (
	internalErrorText = "internal server error"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) HTMLResponse(html []byte) {
	h.rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.rw.WriteHeader(http.StatusOK)

	if _, err := h.rw.Write(html); err != nil {
		h.log.Error("write HTML HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
		sentinel   error
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
		sentinel = core_errors.ErrInvalidArgument

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
		sentinel = core_errors.ErrNotFound

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
		sentinel = core_errors.ErrConflict

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		publicErrorText(err, sentinel),
		msg,
	)
}

func publicErrorText(err error, sentinel error) string {
	current := err
	for {
		next := errors.Unwrap(current)
		if next == nil || next == sentinel {
			return current.Error()
		}

		current = next
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err), zap.Stack("stack"))

	h.errorResponse(
		statusCode,
		err.Error(),
		msg,
	)
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	errText string,
	msg string,
) {
	if statusCode >= http.StatusInternalServerError {
		errText = internalErrorText
	}

	response := ErrorResponse{
		Error:   errText,
		Message: msg,
	}

	h.JSONResponse(
		response,
		statusCode,
	)
}
