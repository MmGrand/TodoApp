package tasks_service

import (
	"context"
	"fmt"

	"github.com/MmGrand/TodoApp/internal/core/domain"
	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	id int,
	version int,
	patch domain.TaskPatch,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if task.Version != version {
		return domain.Task{}, fmt.Errorf(
			"task with id='%d' has version %d, got %d: %w",
			id,
			task.Version,
			version,
			core_errors.ErrConflict,
		)
	}

	if patch.IsEmpty() {
		return task, nil
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}
