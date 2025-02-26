package repositories

import (
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"context"
)

type Basic interface {
	List(ctx context.Context) ([]entities.Basic, error)
	Update(ctx context.Context, basic inputs.UpdateBasic) error
}
