package domain

import (
	"fmt"

	core_errors "github.com/MmGrand/TodoApp/internal/core/errors"
)

const (
	DefaultPaginationLimit = 50
	MaxPaginationLimit     = 200
)

type Pagination struct {
	Limit  int
	Offset int
}

func NewPagination(limit *int, offset *int) (Pagination, error) {
	pagination := Pagination{
		Limit:  DefaultPaginationLimit,
		Offset: 0,
	}

	if limit != nil {
		if *limit < 1 || *limit > MaxPaginationLimit {
			return Pagination{}, fmt.Errorf(
				"limit must be between 1 and %d: %w",
				MaxPaginationLimit,
				core_errors.ErrInvalidArgument,
			)
		}

		pagination.Limit = *limit
	}

	if offset != nil {
		if *offset < 0 {
			return Pagination{}, fmt.Errorf(
				"offset must be non-negative: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		pagination.Offset = *offset
	}

	return pagination, nil
}
