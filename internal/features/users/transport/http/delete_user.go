package users_transport_http

import (
	"net/http"

	core_logger "github.com/nullphosphor/golang-todoapp/internal/core/logger"
	core_http_request "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
