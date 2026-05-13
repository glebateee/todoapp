package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/glebateee/todoapp/internal/core/domain"
	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_request "github.com/glebateee/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/glebateee/todoapp/internal/core/transport/http/response"
	core_http_types "github.com/glebateee/todoapp/internal/core/transport/http/types"
	core_http_utils "github.com/glebateee/todoapp/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'Fullname' can't be null if set")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("'Fullname' length can't be less than 3 and greater than 100")
		}
	}
	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			// 	return fmt.Errorf("'PhoneNumber' can't be null if set")
			// }

			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("'PhoneNumber' length can't be less than 10 and greater than 15")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContextMust(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId path value")
		return
	}

	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	//logger.Debug(fmt.Sprintf("fields:\nFullName: '%s'\nPhoneNumber: '%s'\n", request.FullName, request.PhoneNumber))

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
