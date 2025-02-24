package app

import (
	"checker/internal/config"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func Worker(id int, jobs <-chan config.Server, results chan<- string, logFile string, timeout int) error {
	for server := range jobs {
		err := CheckServer(server.Name, server.Url, logFile, timeout)
		if err != nil {
			results <- fmt.Sprintf("ERROR: [%s] %s", server.Url, err)
		} else {
			results <- fmt.Sprintf("SUCCESS: %s (%s) is healthy", server.Name, server.Url)
		}
	}
	return nil
}

func StartWorkers(cfg *config.Config, workerCount int) error {
	jobs := make(chan config.Server, len(cfg.Servers))
	results := make(chan string, len(cfg.Servers))

	var g errgroup.Group

	for w := 1; w <= workerCount; w++ {
		w := w
		g.Go(func() error {
			return Worker(w, jobs, results, "logs/errors.log", cfg.Basic.Timeout)
		})
	}

	go func() {
		for _, server := range cfg.Servers {
			jobs <- server
		}
		close(jobs)
	}()

	if err := g.Wait(); err != nil {
		return fmt.Errorf("some workers failed: %v", err)
	}

	close(results)

	return nil
}

func TimeScheduler(cfg *config.Config, workerCount int) {
	ticker := time.NewTicker(time.Duration(cfg.Basic.Interval) * time.Second)
	defer ticker.Stop()

	log.Printf("Scheduler running every %d seconds...\n", cfg.Basic.Interval)
	for {
		select {
		case <-ticker.C:
			go func() {
				if err := StartWorkers(cfg, workerCount); err != nil {
					log.Printf("Health check failed: %v", err)
				}
			}()
		}
	}
}
