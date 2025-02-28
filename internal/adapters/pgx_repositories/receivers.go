package pgx_repositories

import (
	app_errors "checker/internal/domain/app/errors"
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxReceiversRepository struct {
	DbPool *pgxpool.Pool
}

func NewPgxReceiversRepository(dbPool *pgxpool.Pool) *PgxReceiversRepository {
	return &PgxReceiversRepository{DbPool: dbPool}
}

func (r *PgxReceiversRepository) Create(ctx context.Context, receivers inputs.CreateReceiver) (int, error) {
	var id int

	query := `INSERT INTO receivers (email) VALUES (@email) RETURNING id;`

	receiverParams := pgx.NamedArgs{
		"email": receivers.Email,
	}

	err := r.DbPool.QueryRow(ctx, query, receiverParams).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PgxReceiversRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM receivers WHERE id = @id;`

	args := pgx.NamedArgs{
		"id": id,
	}

	receivers, err := r.DbPool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if receivers.RowsAffected() == 0 {
		return app_errors.ReceiverDoesNotExist
	}

	return nil
}

func (r *PgxReceiversRepository) List(ctx context.Context) ([]entities.Receiver, error) {
	var receivers []entities.Receiver

	query := `
		SELECT 
			id, 
			email 
		FROM 
			receivers 
		WHERE 
			1 = 1
		`
	rows, err := r.DbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var receiver entities.Receiver
		err := rows.Scan(
			&receiver.ID,
			&receiver.Email,
		)
		if err != nil {
			return nil, err
		}
		receivers = append(receivers, receiver)
	}

	return receivers, nil
}
