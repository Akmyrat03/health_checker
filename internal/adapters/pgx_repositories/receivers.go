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
			email,
			muted
		FROM 
			receivers;`

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
			&receiver.Muted,
		)
		if err != nil {
			return nil, err
		}

		receivers = append(receivers, receiver)
	}

	return receivers, nil
}

func (r *PgxReceiversRepository) GetAllUnmuted(ctx context.Context) ([]entities.Receiver, error) {
	var unmutedReceivers []entities.Receiver

	query := `
		SELECT 
			id, 
			email,
			muted
		FROM 
			receivers
		WHERE
			muted = $1;`

	rows, err := r.DbPool.Query(ctx, query, false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var receiver entities.Receiver
		err := rows.Scan(
			&receiver.ID,
			&receiver.Email,
			&receiver.Muted,
		)
		if err != nil {
			return nil, err
		}

		unmutedReceivers = append(unmutedReceivers, receiver)
	}

	return unmutedReceivers, nil
}

func (r *PgxReceiversRepository) MuteStatus(ctx context.Context, email string, mute bool) error {
	query := `UPDATE receivers SET muted = $1 WHERE email = $2;`
	_, err := r.DbPool.Exec(ctx, query, mute, email)
	if err != nil {
		return err
	}

	return nil
}
