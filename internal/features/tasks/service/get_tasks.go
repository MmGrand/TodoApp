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
) ([]domain.Task, error) {
	pagination, err := domain.NewPagination(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(
		ctx,
		userID,
		&pagination.Limit,
		&pagination.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository: %w", err)
	}

	return tasks, nil
}
