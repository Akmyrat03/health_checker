package repositories

import (
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"context"
)

type Receivers interface {
	Create(ctx context.Context, receivers inputs.CreateReceiver) (int, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]entities.Receiver, error)
}
