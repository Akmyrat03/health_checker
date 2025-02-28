package main

import (
	_ "checker/docs"
	"checker/internal/adapters/pgx_repositories"
	"checker/internal/config"
	"checker/internal/domain/app/usecases/smtp"
	"checker/internal/infrastructure/fiber"
	"checker/internal/infrastructure/pgx"
	"context"
	"log"
)

// @title Health Checker API
// @version 1.0
// @description API for managing servers
// @schemas http, https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name authorization
func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("[config.LoadConfig]: failed to load config file: %v", err)
	}

	// Connect to the database
	pool, err := pgx.PostgresPool()
	if err != nil {
		log.Fatalf("[PostgresPool]: failed to connect to database: %v", err)
	}

	// Initialize repositories
	serverRepo := pgx_repositories.NewPgxRepository(pool)
	basicRepo := pgx_repositories.NewPgxBasicRepository(pool)
	smtpRepo, err := pgx_repositories.NewSMTPRepository(pool, "config.json")
	if err != nil {
		log.Fatalf("[SMTPRepository]: failed to initialize SMTP repository: %v", err)
	}
	scheduler := smtp.Content{
		WorkerCount: 2, // Define worker count
		ServerRepo:  serverRepo,
		BasicRepo:   basicRepo,
		SMTPRepo:    smtpRepo,
	}

	go scheduler.TimeScheduler(context.Background())

	fiber.RunFiberServer(cfg)
}
