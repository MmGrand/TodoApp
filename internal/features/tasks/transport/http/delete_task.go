package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_http_request "github.com/MmGrand/TodoApp/internal/core/transport/http/request"
	core_http_response "github.com/MmGrand/TodoApp/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary Удаление задачи
// @Description Удаление существующей в системе задачи по её ID.
// @Description Если передан `version` и задачу уже изменили, вернётся 409.
// @Tags tasks
// @Param id path int true "ID удаляемой задачи"
// @Param version query int false "Версия задачи, которую видел клиент"
// @Success 204 "Успешное удаление задачи"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict: task was concurrently modified"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'taskID' path value",
		)

		return
	}

	version, err := core_http_request.GetVersionQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'version' query param",
		)

		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID, version); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)

		return
	}

	responseHandler.NoContentResponse()
}
