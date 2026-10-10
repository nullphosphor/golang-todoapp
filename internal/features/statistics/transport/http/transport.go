package statistics_trasnport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
	core_http_server "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from, to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statistics StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statistics,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
