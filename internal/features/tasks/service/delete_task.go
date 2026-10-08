package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	id int,
	version *int,
) error {
	if err := s.tasksRepository.DeleteTask(ctx, id, version); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}

	return nil
}
