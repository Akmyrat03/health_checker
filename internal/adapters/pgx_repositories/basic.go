package pgx_repositories

import (
	"checker/internal/domain/entities"
	"context"

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
		SELECT id, check_interval, timeout, error_interval FROM basic_config WHERE 1 = 1
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
