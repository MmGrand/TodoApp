package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MmGrand/TodoApp/internal/core/domain"
)

func (r *StatisticsRepository) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT
		COUNT(*),
		COUNT(*) FILTER (WHERE completed),
		EXTRACT(EPOCH FROM AVG(completed_at - created_at) FILTER (WHERE completed))::float8
	FROM todoapp.tasks
	`)

	args := []any{}
	conditions := []string{}

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, userID)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	var (
		tasksCreated         int
		tasksCompleted       int
		avgCompletionSeconds *float64
	)

	err := r.pool.QueryRow(ctx, queryBuilder.String(), args...).Scan(
		&tasksCreated,
		&tasksCompleted,
		&avgCompletionSeconds,
	)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("scan statistics: %w", err)
	}

	var tasksCompletedRate *float64
	if tasksCreated > 0 {
		rate := float64(tasksCompleted) / float64(tasksCreated) * 100
		tasksCompletedRate = &rate
	}

	var tasksAverageCompletionTime *time.Duration
	if avgCompletionSeconds != nil {
		avg := time.Duration(*avgCompletionSeconds * float64(time.Second)).Round(time.Millisecond)
		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		tasksCompletedRate,
		tasksAverageCompletionTime,
	), nil
}
