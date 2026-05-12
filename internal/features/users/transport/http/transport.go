package users_transport_http

import (
	"context"
	"net/http"

	"github.com/glebateee/todoapp/internal/core/domain"
	core_http_server "github.com/glebateee/todoapp/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService usersService
}

type usersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)
}

func NewUsersHTTPHandler(
	usersService usersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
	}
}
