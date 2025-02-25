package usecases

import (
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/app/repositories"
	"checker/internal/domain/entities"
	"context"
	"fmt"
)

type ServerUseCase struct {
	serverRepository repositories.Server
}

func NewServerUseCase(serverRepository repositories.Server) *ServerUseCase {
	return &ServerUseCase{serverRepository: serverRepository}
}

func (serverUsecase *ServerUseCase) Create(ctx context.Context, server inputs.CreateServer) (int, error) {
	id, err := serverUsecase.serverRepository.Create(ctx, server)
	if err != nil {
		return 0, fmt.Errorf("failed to create server: %v", err)
	}

	return id, nil
}

func (serverUseCase *ServerUseCase) Delete(ctx context.Context, id int) error {
	err := serverUseCase.serverRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete server by id")
	}

	return nil
}

func (serverUseCase *ServerUseCase) List(ctx context.Context) ([]entities.Server, error) {
	servers, err := serverUseCase.serverRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return servers, nil
}
