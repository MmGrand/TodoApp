package core_pgx_pool

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/MmGrand/TodoApp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

func (r pgxRows) Scan(dest ...any) error {
	if err := r.Rows.Scan(dest...); err != nil {
		return mapErrors(err)
	}

	return nil
}

func (r pgxRows) Err() error {
	if err := r.Rows.Err(); err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
		pgxViolatesCheckErrorCode      = "23514"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgxViolatesForeignKeyErrorCode:
			return fmt.Errorf("%w %q", core_postgres_pool.ErrViolatesForeignKey, pgErr.ConstraintName)

		case pgxViolatesCheckErrorCode:
			return fmt.Errorf("%w %q", core_postgres_pool.ErrViolatesCheck, pgErr.ConstraintName)
		}
	}

	return fmt.Errorf("%w: %w", err, core_postgres_pool.ErrUnknown)
}
