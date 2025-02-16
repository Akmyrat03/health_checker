package health

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

func worker(id int, jobs <-chan string, results chan<- string, logFile string, timeout int, g *errgroup.Group) {
	g.Go(func() error {
		for server := range jobs {
			err := CheckServerHealth(server, logFile, timeout)
			if err != nil {
				results <- fmt.Sprintf("ERROR: %s", err)
			} else {
				results <- fmt.Sprintf("SUCCESS: %s is healthy", server)
			}
		}
		return nil
	})
}

func StartWorkers(servers []string, workerCount int, checkInterval int, logFile string, timeout int) error {
	jobs := make(chan string, len(servers))
	results := make(chan string, len(servers))

	var g errgroup.Group

	for w := 1; w <= workerCount; w++ {
		go worker(w, jobs, results, logFile, timeout, &g)
	}

	for _, server := range servers {
		jobs <- server
	}
	close(jobs)

	if err := g.Wait(); err != nil {
		return fmt.Errorf("some workers failed: %v", err)
	}

	// close(results)

	for range servers {
		<-results
	}
	close(results)

	return nil
}
