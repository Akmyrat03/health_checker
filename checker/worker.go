package health

import (
	"checker/config"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func worker(id int, jobs <-chan string, results chan<- string, logFile string, timeout int) error {
	for server := range jobs {
		err := CheckServerHealth(server, logFile, timeout)
		if err != nil {
			results <- fmt.Sprintf("ERROR: %s", err)
		} else {
			results <- fmt.Sprintf("SUCCESS: %s is healthy", server)
		}
	}
	return nil
}

func StartWorkers(filename string, workerCount int, checkInterval int, logFile string, timeout int) error {
	jobs := make(chan string, workerCount)
	results := make(chan string, workerCount)

	var g errgroup.Group

	for w := 1; w <= workerCount; w++ {
		w := w
		g.Go(func() error {
			return worker(w, jobs, results, logFile, timeout)
		})
	}

	serversChan, err := config.LoadServers(filename)
	if err != nil {
		return fmt.Errorf("error loading servers: %v", err)
	}

	go func() {
		for server := range serversChan {
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
