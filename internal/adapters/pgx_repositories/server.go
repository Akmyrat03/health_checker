package pgx_repositories

import (
	app_errors "checker/internal/domain/app/errors"
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxRepository struct {
	DbPool *pgxpool.Pool
}

func NewPgxRepository(dbPool *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{DbPool: dbPool}
}

func (r *PgxRepository) Create(ctx context.Context, server inputs.CreateServer) (int, error) {
	var id int

	addServerQuery := `
		INSERT INTO servers (
			name,
			url
		) VALUES (
			@name,
			@url)
		RETURNING id;
	`

	serverParams := pgx.NamedArgs{
		"name": server.Name,
		"url":  server.Url,
	}

	err := r.DbPool.QueryRow(ctx, addServerQuery, serverParams).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PgxRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM servers WHERE id = @id;
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	servers, err := r.DbPool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if servers.RowsAffected() == 0 {
		return app_errors.ServerDoesNotExist
	}

	return nil
}

func (r *PgxRepository) List(ctx context.Context) ([]entities.Server, error) {
	var servers []entities.Server

	query := `
		SELECT id, name, url FROM servers WHERE 1 = 1
	`

	rows, err := r.DbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var server entities.Server
		err := rows.Scan(
			&server.ID,
			&server.Name,
			&server.Url,
		)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}

	return servers, nil
}
