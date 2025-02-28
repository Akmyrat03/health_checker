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

func (basicUseCase *BasicUseCase) Get(ctx context.Context) (*entities.Basic, error) {
	basic, err := basicUseCase.basicRepository.Get(ctx)
	if err != nil {
		return nil, err
	}

	return basic, nil
}

func (basicUseCase *BasicUseCase) Update(ctx context.Context, basic inputs.UpdateBasic) error {
	err := basicUseCase.basicRepository.Update(ctx, basic)
	if err != nil {
		return err
	}

	return nil
}
