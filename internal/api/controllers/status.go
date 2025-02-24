package controllers

import (
	"checker/internal/adapters/smtp"
	"checker/internal/config"
	"checker/internal/domain/app/usecases"

	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type HealthStatus struct {
	ServerName string `json:"server_name"`
	ServerURL  string `json:"server_url"`
	Status     string `json:"status"`
}

func ShowStatus(c *fiber.Ctx) error {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("[status.go]: error loading config file: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	messageSender, err := smtp.NewSMTP("config.json")
	if err != nil {
		fmt.Printf("Failed to initialize SMTP: %v\n", err)
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	status := []HealthStatus{}

	for _, server := range cfg.Servers {
		wg.Add(1)

		go func(server config.Server) {
			defer wg.Done()

			err := usecases.CheckServer(server.Name, server.Url, "logs/errors.log", cfg.Basic.Timeout, messageSender)
			serverStatus := "Healthy"
			if err != nil {
				serverStatus = "Unhealthy"
			}

			mu.Lock()
			status = append(status, HealthStatus{
				ServerName: server.Name,
				ServerURL:  server.Url,
				Status:     serverStatus,
			})
			mu.Unlock()
		}(server)
	}

	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"servers": status,
	})
}
