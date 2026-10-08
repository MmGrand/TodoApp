package users_postgres_repository

import (
	"context"
	"fmt"

	core_postgres_pool "github.com/MmGrand/TodoApp/internal/core/repository/postgres/pool"
)

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}

func (r *UsersRepository) userExists(ctx context.Context, id int) (bool, error) {
	query := `
	SELECT EXISTS(SELECT 1 FROM todoapp.users WHERE id = $1);
	`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("check user existence: %w", err)
	}

	return exists, nil
}
