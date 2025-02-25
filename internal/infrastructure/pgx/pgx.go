package pgx

import (
	"checker/internal/config"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	once        sync.Once
	connectPool *pgxpool.Pool
)

func PostgresPool() (*pgxpool.Pool, error) {
	once.Do(func() {
		cfg, err := config.LoadConfig("config.json")
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
		dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.DBName,
			cfg.Postgres.SslMode,
		)

		connectConfig, err := pgxpool.ParseConfig(dbUrl)
		if err != nil {
			fmt.Printf("Unable to parse postgres connection string: %s\n", dbUrl)
		}

		connectPool, err = pgxpool.NewWithConfig(context.Background(), connectConfig)
		if err != nil {
			fmt.Printf("Unable to connect to database: %s\n", dbUrl)
		}
	})

	return connectPool, nil
}
