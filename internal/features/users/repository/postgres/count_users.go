package users_postgres_repository

import (
	"context"
	"fmt"
)

func (r *UsersRepository) CountUsers(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT COUNT(*)
	FROM todoapp.users;
	`

	var count int
	if err := r.pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}
