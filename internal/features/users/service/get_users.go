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
) ([]domain.User, error) {
	pagination, err := domain.NewPagination(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	users, err := s.usersRepository.GetUsers(
		ctx,
		&pagination.Limit,
		&pagination.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil
}
