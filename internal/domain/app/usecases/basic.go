package usecases

import (
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/app/repositories"
	"checker/internal/domain/entities"
	"context"
)

type BasicUseCase struct {
	basicRepository repositories.Basic
}

func NewBasicUseCase(basicRepository repositories.Basic) *BasicUseCase {
	return &BasicUseCase{basicRepository: basicRepository}
}

func (basicUseCase *BasicUseCase) List(ctx context.Context) ([]entities.Basic, error) {
	basics, err := basicUseCase.basicRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return basics, nil
}

func (basicUseCase *BasicUseCase) Update(ctx context.Context, basic inputs.UpdateBasic) error {
	err := basicUseCase.basicRepository.Update(ctx, basic)
	if err != nil {
		return err
	}

	return nil
}
