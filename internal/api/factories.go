package api

import (
	"checker/internal/adapters/pgx_repositories"
	"checker/internal/api/providers"
	"checker/internal/domain/app/usecases"
	"fmt"
)

func MakeServerUseCase() (*usecases.ServerUseCase, error) {
	pool, err := providers.GetDbPool()
	if err != nil {
		fmt.Printf("ERROR: Failed to get database pool: %v\n", err)
		return nil, err
	}
	repo := pgx_repositories.NewPgxRepository(pool)
	serverUseCase := usecases.NewServerUseCase(repo)
	return serverUseCase, nil
}

func MakeBasicUseCase() (*usecases.BasicUseCase, error) {
	pool, err := providers.GetDbPool()
	if err != nil {
		fmt.Printf("ERROR: Failed to get dababase pool: %v\n", err)
		return nil, err
	}
	repo := pgx_repositories.NewPgxBasicRepository(pool)
	basicUseCase := usecases.NewBasicUseCase(repo)
	return basicUseCase, nil
}

func MakeReceiverUseCase() (*usecases.ReceiversUseCase, error) {
	pool, err := providers.GetDbPool()
	if err != nil {
		fmt.Printf("ERROR: Failed to get database pool: %v\n", err)
		return nil, err
	}
	repo := pgx_repositories.NewPgxReceiversRepository(pool)
	receiverUseCase := usecases.NewReceiversUseCase(repo)
	return receiverUseCase, nil
}

// func MakeSMTPUseCase() (*usecases.SMTPUseCase, error) {
// 	pool, err := providers.GetDbPool()
// 	if err != nil {
// 		fmt.Printf("ERROR: Failed to get database pool: %v\n", err)
// 		return nil, err
// 	}
// 	repo, err := pgx_repositories.NewSMTPRepository(pool, "config.json")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return usecases.NewSMTPUseCase(repo), nil

// }
