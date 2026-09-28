package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	if err := s.tasksRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}

	return nil
}
