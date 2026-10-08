package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/MmGrand/TodoApp/internal/core/domain"
	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_http_request "github.com/MmGrand/TodoApp/internal/core/transport/http/request"
	core_http_response "github.com/MmGrand/TodoApp/internal/core/transport/http/response"
	core_http_types "github.com/MmGrand/TodoApp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Version     int                              `json:"version" example:"3"`
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Ivan Petrov"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887766"`
}

func (r *PatchUserRequest) Validate() error {
	if r.Version < 1 {
		return fmt.Errorf("`version` is required and must be positive")
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary Изменение пользователя
// @Description Обновить поля существующего пользователя по его ID.
// @Description Передаются только изменяемые поля: отсутствующее поле не меняется,
// @Description `null` в `phone_number` очищает номер телефона. `full_name` не может быть `null`.
// @Description Обязательное поле `version` — версия пользователя, которую видел клиент; если его уже изменили, вернётся 409.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID изменяемого пользователя"
// @Param request body PatchUserRequest true "PatchUser тело запроса"
// @Success 200 {object} PatchUserResponse "Обновлённый пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict: user was concurrently modified"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'userID' path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, request.Version, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
