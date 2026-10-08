package users_service

import (
	"context"
	"fmt"

	"github.com/MmGrand/TodoApp/internal/core/domain"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, int, error) {
	pagination, err := domain.NewPagination(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("pagination: %w", err)
	}

	users, err := s.usersRepository.GetUsers(
		ctx,
		&pagination.Limit,
		&pagination.Offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("get users from repository: %w", err)
	}

	total, err := s.usersRepository.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count users in repository: %w", err)
	}

	return users, total, nil
}
