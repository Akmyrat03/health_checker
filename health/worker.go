package health

import "fmt"

func worker(id int, jobs <-chan string, results chan<- string, logFile string) {
	for server := range jobs {
		err := CheckServerHealth(server, logFile)
		if err != nil {
			results <- fmt.Sprintf("ERROR: %s", err)
		} else {
			results <- fmt.Sprintf("SUCCESS: %s is healthy", server)
		}
	}
}

func StartWorkers(servers []string, workerCount int, checkInterval int, logFile string) {
	jobs := make(chan string, len(servers))
	results := make(chan string, len(servers))

	for w := 1; w <= workerCount; w++ {
		go worker(w, jobs, results, logFile)
	}

	for _, server := range servers {
		jobs <- server
	}
	close(jobs)

	for range servers {
		<-results
	}
}
