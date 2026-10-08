package tasks_postgres_repository

import (
	"context"
	"fmt"
	"strings"
)

func (r *TasksRepository) CountTasks(
	ctx context.Context,
	userID *int,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT COUNT(*)
	FROM todoapp.tasks
	`)

	args := []any{}

	if userID != nil {
		args = append(args, userID)
		fmt.Fprintf(&queryBuilder, " WHERE author_user_id=$%d", len(args))
	}

	var count int
	if err := r.pool.QueryRow(ctx, queryBuilder.String(), args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}

	return count, nil
}
