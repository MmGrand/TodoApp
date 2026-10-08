package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_postgres_pool "github.com/MmGrand/TodoApp/internal/core/repository/postgres/pool"
)

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}

func (r *TasksRepository) taskExists(ctx context.Context, id int) (bool, error) {
	query := `
	SELECT EXISTS(SELECT 1 FROM todoapp.tasks WHERE id = $1);
	`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("check task existence: %w", err)
	}

	return exists, nil
}
