package fiber

import (
	"checker/internal/adapters/pgx_repositories"
	"checker/internal/api"
	rest_v0 "checker/internal/api/rest/v0"
	"checker/internal/config"
	"checker/internal/infrastructure/pgx"
	"checker/internal/infrastructure/scheduler"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func RunFiberServer(cfg *config.Config) {
	var err error

	pool, err := pgx.PostgresPool()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	basicRepo := pgx_repositories.NewPgxBasicRepository(pool)
	serverRepo := pgx_repositories.NewPgxRepository(pool)

	basicConfig, err := basicRepo.Get(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch basic config, skipped details: %v", err)
	} else {
		fmt.Printf("Server started on : %s:%s\n", cfg.App.Host, cfg.App.Port)
		fmt.Printf("Check Interval: %d seconds\n", basicConfig.CheckInterval)
		fmt.Printf("Timeout: %d seconds\n", basicConfig.Timeout)
		fmt.Printf("Notification Interval: %d hours\n", basicConfig.NotificationInterval)
	}

	receiverUseCase, err := api.MakeReceiverUseCase()
	if err != nil {
		fmt.Printf("[fiber.go]: api.MakeReceiverUseCase : %v", err)
		return
	}

	scheduler := scheduler.Content{
		WorkerCount:     2, // Define worker count
		ServerRepo:      serverRepo,
		BasicRepo:       basicRepo,
		ReceiverUseCase: receiverUseCase,
	}

	go scheduler.TimeScheduler(context.Background())

	app := fiber.New(
		fiber.Config{
			BodyLimit: 50 * 1024 * 1024,
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Health Checker Service!")
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Cors.Origins,
		AllowCredentials: cfg.Cors.Credentials,
		AllowHeaders:     "*",
		AllowMethods:     "GET POST PUT DELETE",
	}))

	api := app.Group("/api")
	rest_v0.GroupControllers(&api)

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		if pool == nil {
			return fmt.Errorf("database pool is not initialized")
		}

		return nil
	})

	app.Hooks().OnShutdown(func() error {
		if pool != nil {
			pool.Close()
		}

		return nil
	})

	go func() {
		err = app.Listen(
			fmt.Sprintf("%s:%s",
				cfg.App.Host,
				cfg.App.Port,
			))
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
