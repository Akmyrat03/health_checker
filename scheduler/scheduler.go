package scheduler

import (
	"checker/health"
	"log"
	"time"
)

func StartHealthCheckScheduler(servers []string, checkInterval int, logFile string, timeout int) {
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			if err := health.StartWorkers(servers, 2, checkInterval, logFile, timeout); err != nil {
				log.Printf("Health check failed: %v", err)
			}
		}
	}
}
