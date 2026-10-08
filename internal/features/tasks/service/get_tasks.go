package tasks_service

import (
	"context"
	"fmt"

	"github.com/MmGrand/TodoApp/internal/core/domain"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, int, error) {
	pagination, err := domain.NewPagination(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("pagination: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(
		ctx,
		userID,
		&pagination.Limit,
		&pagination.Offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("get tasks from repository: %w", err)
	}

	total, err := s.tasksRepository.CountTasks(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("count tasks in repository: %w", err)
	}

	return tasks, total, nil
}
