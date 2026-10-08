package statistics_transport_http

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/MmGrand/TodoApp/internal/core/domain"
	core_logger "github.com/MmGrand/TodoApp/internal/core/logger"
	core_http_request "github.com/MmGrand/TodoApp/internal/core/transport/http/request"
	core_http_response "github.com/MmGrand/TodoApp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"20"`
	TasksCompleted             int      `json:"tasks_completed" example:"15"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"75.25"`
	TasksAverageCompletionTime *float64 `json:"tasks_average_completion_seconds" example:"95400.5"`
}

// GetStatistics godoc
// @Summary Статистика по задачам
// @Description Получить статистику по задачам: количество созданных и выполненных,
// @Description процент выполнения (округлён до сотых) и среднее время выполнения в секундах.
// @Description Можно отфильтровать по автору и периоду создания задач `[from, to)`; `to` должен быть позже `from`.
// @Description Даты трактуются в часовом поясе приложения (TIME_ZONE).
// @Tags statistics
// @Produce json
// @Param user_id query int false "ID автора задач"
// @Param from query string false "Начало периода (YYYY-MM-DD)" Format(date)
// @Param to query string false "Конец периода (YYYY-MM-DD), не включается: задачи, созданные строго до этой даты" Format(date)
// @Success 200 {object} GetStatisticsResponse "Статистика по задачам"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r, h.location)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'userID'/'from'/'to' query param",
		)

		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)

		return
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var rate *float64
	if statistics.TasksCompletedRate != nil {
		rounded := math.Round(*statistics.TasksCompletedRate*100) / 100
		rate = &rounded
	}

	var avgSeconds *float64
	if statistics.TasksAverageCompletionTime != nil {
		seconds := statistics.TasksAverageCompletionTime.Seconds()
		avgSeconds = &seconds
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         rate,
		TasksAverageCompletionTime: avgSeconds,
	}
}

func getUserIDFromToQueryParams(
	r *http.Request,
	location *time.Location,
) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey, location)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey, location)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
