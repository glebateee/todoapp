package users_transport_http

import (
	"net/http"

	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_response "github.com/glebateee/todoapp/internal/core/transport/http/response"
)

type GetUserresponse UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContextMust(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)
}
