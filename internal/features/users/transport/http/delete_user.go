package users_transport_http

import (
	"net/http"

	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_request "github.com/glebateee/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/glebateee/todoapp/internal/core/transport/http/response"
)

// DELETE /users/{id}

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContextMust(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	userId, err := core_http_request.GetIntPathValue(r, "id")
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
