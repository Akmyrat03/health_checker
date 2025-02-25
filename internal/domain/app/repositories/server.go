package repositories

import (
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"context"
)

type Server interface {
	Create(ctx context.Context, server inputs.CreateServer) (int, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]entities.Server, error)
}
