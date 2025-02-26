package main

import (
	_ "checker/docs"
	"checker/internal/adapters/smtp"
	"checker/internal/config"
	"checker/internal/domain/app/usecases"
	"checker/internal/infrastructure/fiber"
	"fmt"
	"log"
)

// @title Health Checker API
// @version 1.0
// @description API for managing servers
// @schemas http, https
// @in header
func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("[config.LoadConfig]: failed to load config file: %v", err)
	}

	messageSender, err := smtp.NewSMTP("config.json")
	if err != nil {
		fmt.Printf("Failed to initialize SMTP: %v\n", err)
		return
	}

	go usecases.TimeScheduler(cfg, 2, messageSender)

	fiber.RunFiberServer(cfg)
}
