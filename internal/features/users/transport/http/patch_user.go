package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
	core_logger "github.com/nullphosphor/golang-todoapp/internal/core/logger"
	core_http_request "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Fullname    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (p *PatchUserRequest) Validate() error {
	if p.Fullname.Set {
		if p.Fullname.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		if fullNameLen := len([]rune(*p.Fullname.Value)); fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if p.PhoneNumber.Set && p.PhoneNumber.Value != nil {
		if phoneNumberLen := len([]rune(*p.PhoneNumber.Value)); phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
		}

		if !strings.HasPrefix(*p.PhoneNumber.Value, "+") {
			return fmt.Errorf("`PhoneNumber` must start with symbol '+'")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
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

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP PATCH request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(UserDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Fullname.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
