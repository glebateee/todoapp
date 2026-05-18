package tasks_transport_http

import (
	"net/http"

	"github.com/glebateee/todoapp/internal/core/domain"
	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_request "github.com/glebateee/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/glebateee/todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorUserId int     `json:"author_user_id" validate:"required"`
}

type CreateTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContextMust(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	var request CreateTaskRequest

	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserId,
	)

	domainTask, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}
	response := CreateTaskResponse(taskDTOFromDomain(domainTask))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
