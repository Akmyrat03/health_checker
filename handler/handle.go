package handler

import (
	"checker/config"
	"checker/health"
	"encoding/json"
	"net/http"
)

type HealthStatus struct {
	ServerURL string `json:"server_url"`
	Status    string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request, config *config.Config) {
	status := make([]HealthStatus, 0)

	for _, server := range config.Servers {
		err := health.CheckServerHealth(server, config.LogFile, config.Timeout)
		status = append(status, HealthStatus{
			ServerURL: server,
			Status:    "Healthy",
		})

		if err != nil {
			status[len(status)-1].Status = "Unhealthy"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
