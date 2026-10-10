package statistics_service

import (
	"context"
	"time"

	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from, to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
