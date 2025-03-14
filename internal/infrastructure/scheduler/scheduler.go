package scheduler

import (
	"checker/internal/domain/app/repositories"
	"checker/internal/domain/app/usecases"
	"checker/internal/domain/entities"
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type Content struct {
	WorkerCount     int
	ServerRepo      repositories.Server
	BasicRepo       repositories.Basic
	ReceiverUseCase *usecases.ReceiversUseCase
}

func (c *Content) Worker(ctx context.Context, jobs <-chan entities.Server, results chan<- string) error {
	for server := range jobs {
		// log.Printf("Worker processing server: %s (%s)", server.Name, server.Url)
		err := CheckServer(ctx, server, c.BasicRepo, c.ReceiverUseCase)
		if err != nil {
			results <- fmt.Sprintf("ERROR: [%s] %s", server.Url, err)
		} else {
			results <- fmt.Sprintf("SUCCESS: %s (%s) is healthy", server.Name, server.Url)
		}
	}
	return nil
}

func (c *Content) StartWorkers(ctx context.Context) error {
	servers, err := c.ServerRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %v", err)
	}

	jobs := make(chan entities.Server, len(servers))
	results := make(chan string, len(servers))

	var g errgroup.Group

	for w := 1; w <= c.WorkerCount; w++ {
		g.Go(func() error {
			return c.Worker(ctx, jobs, results)
		})
	}

	go func() {
		for _, server := range servers {
			// log.Printf("Dispatching server check job for: %s (%s)", server.Name, server.Url)
			jobs <- server
		}
		close(jobs)
	}()

	if err := g.Wait(); err != nil {
		return fmt.Errorf("some workers failed: %v", err)
	}

	close(results)

	// for result := range results {
	// 	log.Println("Result:", result)
	// }

	return nil
}

func (c *Content) TimeScheduler(ctx context.Context) {
	basic, err := c.BasicRepo.Get(ctx)
	if err != nil {
		fmt.Printf("Failed to get basic config: %v", err)
		return
	}

	interval := time.Duration(basic.CheckInterval) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				if err := c.StartWorkers(ctx); err != nil {
					log.Printf("Health check failed: %v", err)
				}
			}()
		}
	}
}
