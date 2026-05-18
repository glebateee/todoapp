package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/glebateee/todoapp/internal/core/domain"
	core_logger "github.com/glebateee/todoapp/internal/core/logger"
	core_http_request "github.com/glebateee/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/glebateee/todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int `json:""`
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *string
}

func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContextMust(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID/from/to query params")
		return
	}
	domainStatistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}
	response := toDTOFromDomain(domainStatistics)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		d := statistics.TasksAverageCompletionTime.String()
		avgTime = &d
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	userId, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}
	return userId, from, to, nil
}
