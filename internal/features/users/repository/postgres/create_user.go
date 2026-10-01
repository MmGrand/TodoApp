package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MmGrand/TodoApp/internal/core/domain"
	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
	core_postgres_pool "github.com/MmGrand/TodoApp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesCheck) {
			return domain.User{}, fmt.Errorf("%w: %w", core_postgres_pool.ErrViolatesCheck, core_errors.ErrInvalidArgument)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
