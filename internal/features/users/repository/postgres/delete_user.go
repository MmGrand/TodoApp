package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
	version *int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id = $1 AND ($2::bigint IS NULL OR version = $2);
	`

	cmdTag, err := r.pool.Exec(ctx, query, id, version)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() > 0 {
		return nil
	}

	if version != nil {
		exists, err := r.userExists(ctx, id)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf(
				"user with id='%d' concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}
	}

	return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
}
