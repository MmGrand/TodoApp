package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/MmGrand/TodoApp/internal/core/domain"
	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_http_request "github.com/MmGrand/TodoApp/internal/core/transport/http/request"
	core_http_response "github.com/MmGrand/TodoApp/internal/core/transport/http/response"
	core_http_types "github.com/MmGrand/TodoApp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Version     int                              `json:"version" example:"2"`
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Купить продукты и воду"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Молоко, хлеб, яйца, вода"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"true"`
}

// Validate проверяет только транспортные поля,
// значения задачи валидирует домен при применении патча.
func (r *PatchTaskRequest) Validate() error {
	if r.Version < 1 {
		return fmt.Errorf("`version` is required and must be positive")
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask godoc
// @Summary Изменение задачи
// @Description Обновить поля существующей задачи по её ID.
// @Description Передаются только изменяемые поля: отсутствующее поле не меняется,
// @Description `null` в `description` очищает описание. `title` и `completed` не могут быть `null`.
// @Description Обязательное поле `version` — версия задачи, которую видел клиент; если задачу уже изменили, вернётся 409.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID изменяемой задачи"
// @Param request body PatchTaskRequest true "PatchTask тело запроса"
// @Success 200 {object} PatchTaskResponse "Обновлённая задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict: task was concurrently modified"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, request.Version, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
