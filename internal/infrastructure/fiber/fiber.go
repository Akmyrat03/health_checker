package fiber

import (
	controllers "checker/internal/api"
	"checker/internal/config"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RunFiberServer(cfg *config.Config) {
	var err error

	fmt.Println("Server started on :3000")
	fmt.Println("Check Interval:", cfg.Basic.Interval)
	fmt.Println("Timeout:", cfg.Basic.Timeout)

	app := fiber.New(
		fiber.Config{
			BodyLimit: 50 * 1024 * 1024,
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Post Management Microservice!")
	})

	api := app.Group("/api")
	controllers.GroupControllers(&api)

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

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
