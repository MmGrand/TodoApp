package users_transport_http

import (
	"net/http"
	"strconv"

	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_http_request "github.com/MmGrand/TodoApp/internal/core/transport/http/request"
	core_http_response "github.com/MmGrand/TodoApp/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers godoc
// @Summary Список пользователей
// @Description Получить список пользователей системы с поддержкой пагинации
// @Tags users
// @Produce json
// @Param limit query int false "Максимальное количество пользователей в ответе (1–200, по умолчанию 50)"
// @Param offset query int false "Количество пропускаемых пользователей"
// @Success 200 {object} GetUsersResponse "Список пользователей"
// @Header 200 {integer} X-Total-Count "Общее количество пользователей"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset' query param",
		)

		return
	}

	userDomains, total, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)

		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	rw.Header().Set(core_http_response.TotalCountHeader, strconv.Itoa(total))

	responseHandler.JSONResponse(response, http.StatusOK)
}
