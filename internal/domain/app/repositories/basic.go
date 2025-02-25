package repositories

import (
	"checker/internal/domain/entities"
	"context"
)

type Basic interface {
	List(ctx context.Context) ([]entities.Basic, error)
}
