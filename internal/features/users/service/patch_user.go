package users_service

import (
	"context"
	"fmt"

	"github.com/MmGrand/TodoApp/internal/core/domain"
	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id int,
	version int,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	if user.Version != version {
		return domain.User{}, fmt.Errorf(
			"user with id='%d' has version %d, got %d: %w",
			id,
			user.Version,
			version,
			core_errors.ErrConflict,
		)
	}

	if patch.IsEmpty() {
		return user, nil
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	return patchedUser, nil
}
