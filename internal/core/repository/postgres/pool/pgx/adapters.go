package core_pgx_pool

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/glebateee/todoapp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(args ...any) error {

	err := r.Row.Scan(args...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

type pgxCommandtag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgxViolatesForeignKeyErrorCode {
		return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrViolatesForeignKey)
	}

	return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)
}
