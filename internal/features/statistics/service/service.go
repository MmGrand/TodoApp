package statistics_service

import (
	"context"
	"time"

	"github.com/MmGrand/TodoApp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
