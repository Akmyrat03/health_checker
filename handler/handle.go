package handler

import (
	health "checker/checker"
	"checker/config"
	"encoding/json"
	"net/http"
	"sync"
)

type HealthStatus struct {
	ServerURL string `json:"server_url"`
	Status    string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	serversChan, err := config.LoadServers("servers.json")
	if err != nil {
		http.Error(w, "failed to load servers", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	var status []HealthStatus

	for server := range serversChan {
		wg.Add(1)

		go func(server string) {
			defer wg.Done()

			err := health.CheckServerHealth(server, cfg.LogFile, cfg.Timeout)
			serverStatus := "Healthy"
			if err != nil {
				serverStatus = "Unhealthy"
			}

			status = append(status, HealthStatus{
				ServerURL: server,
				Status:    serverStatus,
			})
		}(server)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
