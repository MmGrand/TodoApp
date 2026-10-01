package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/MmGrand/TodoApp/internal/core/domain"
	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	statistics, err := s.statisticsRepository.GetStatistics(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get statistics from repository: %w", err)
	}

	return statistics, nil
}
