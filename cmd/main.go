package main

import (
	"checker/config"
	"checker/handler"
	"checker/scheduler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("config.LoadConfig - Error: %v", err)
	}

	fmt.Println("Servers:", config.Servers)
	fmt.Println("Check Interval:", config.CheckInterval)
	fmt.Println("Timeout:", config.Timeout)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handler.HealthHandler(w, r, config)
	})

	go http.ListenAndServe(":8080", nil)

	scheduler.StartHealthCheckScheduler(config.Servers, config.CheckInterval, config.LogFile, config.Timeout)
}
