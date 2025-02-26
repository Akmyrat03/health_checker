package pgx_repositories

import (
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxBasicRepository struct {
	DbPool *pgxpool.Pool
}

func NewPgxBasicRepository(dbPool *pgxpool.Pool) *PgxBasicRepository {
	return &PgxBasicRepository{DbPool: dbPool}
}

func (r *PgxBasicRepository) List(ctx context.Context) ([]entities.Basic, error) {
	var basics []entities.Basic

	query := `
		SELECT check_interval, timeout, error_interval FROM basic_config WHERE 1 = 1
	`

	rows, err := r.DbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var basic entities.Basic
		err := rows.Scan(
			&basic.CheckInterval,
			&basic.Timeout,
			&basic.ErrorInterval,
		)
		if err != nil {
			return nil, err
		}
		basics = append(basics, basic)
	}

	return basics, nil
}

func (r *PgxBasicRepository) Update(ctx context.Context, basic inputs.UpdateBasic) error {
	query := `UPDATE basic_config SET check_interval = @check_interval, timeout = @timeout, error_interval = @error_interval`

	args := pgx.NamedArgs{
		"check_interval": basic.CheckInterval,
		"timeout":        basic.Timeout,
		"error_interval": basic.NotificationInterval,
	}

	_, err := r.DbPool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update basic config: %w", err)
	}

	return nil
}
