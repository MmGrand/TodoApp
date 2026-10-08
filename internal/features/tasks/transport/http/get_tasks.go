package tasks_transport_http

import (
	"fmt"
	"net/http"
	"strconv"

	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_http_request "github.com/MmGrand/TodoApp/internal/core/transport/http/request"
	core_http_response "github.com/MmGrand/TodoApp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Список задач
// @Description Получить список задач с поддержкой пагинации и фильтрацией по автору
// @Tags tasks
// @Produce json
// @Param user_id query int false "ID автора задач"
// @Param limit query int false "Максимальное количество задач в ответе (1–200, по умолчанию 50)"
// @Param offset query int false "Количество пропускаемых задач"
// @Success 200 {object} GetTasksResponse "Список задач"
// @Header 200 {integer} X-Total-Count "Общее количество задач с учётом фильтра"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'userID'/'limit'/'offset' query param",
		)

		return
	}

	tasksDomains, total, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))

	rw.Header().Set(core_http_response.TotalCountHeader, strconv.Itoa(total))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const userIDQueryParamKey = "user_id"

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		return nil, nil, nil, err
	}

	return userID, limit, offset, nil
}
