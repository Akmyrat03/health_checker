package providers

import (
	"checker/internal/infrastructure/pgx"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDbPool() (*pgxpool.Pool, error) {
	pool, err := pgx.PostgresPool()
	if err != nil {
		return nil, err
	}

	return pool, err
}
