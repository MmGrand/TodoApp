package tasks_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/MmGrand/TodoApp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	`)

	args := []any{}

	if userID != nil {
		args = append(args, userID)
		fmt.Fprintf(&queryBuilder, " WHERE author_user_id=$%d", len(args))
	}

	args = append(args, limit, offset)
	fmt.Fprintf(&queryBuilder, " ORDER BY id ASC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.pool.Query(
		ctx,
		queryBuilder.String(),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)

		if err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
