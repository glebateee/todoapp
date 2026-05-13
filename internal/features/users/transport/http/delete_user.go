package users_transport_http

import (
	"net/http"

	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_response "github.com/glebateee/todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/glebateee/todoapp/internal/core/transport/http/utils"
)

// DELETE /users/{id}

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContextMust(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userId path value")
		return
	}

	if err := h.usersService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.NoContentResponse()
}
