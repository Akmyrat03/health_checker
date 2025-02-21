package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	limit = 10 * 1024 * 1024 // File size limit for large server files
)

// Config represents the structure of the config.json file
type Config struct {
	HealthChecker struct {
		CheckInterval int    `json:"check_interval"`
		LogFile       string `json:"log_file"`
		Timeout       int    `json:"timeout"`
	} `json:"health_checker"`
}

type Server struct {
	Name   string `json:"name"`
	Server string `json:"server"`
}

// LoadConfig reads the config from the given JSON file and returns the relevant fields
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error while decoding config file: %v", err)
	}

	return &config, nil
}

// LoadServers returns a channel of servers to be processed
func LoadServers(filename string) (<-chan Server, error) {
	serversChan := make(chan Server)

	go func() {
		defer close(serversChan)

		info, err := os.Stat(filename)
		if err != nil {
			fmt.Println("error getting file info:", err)
			return
		}

		// If the file is large, load line by line, otherwise load the full file
		if info.Size() > limit {
			err := LoadServersLineByLine(filename, serversChan)
			if err != nil {
				fmt.Println("error loading servers line by line:", err)
			}
		} else {
			servers, err := LoadServersFull(filename)
			if err != nil {
				fmt.Println("error loading servers fully:", err)
			}
			for _, server := range servers {
				serversChan <- server
			}
		}
	}()

	return serversChan, nil
}

// LoadServersLineByLine loads server data from a large JSON file line by line
func LoadServersLineByLine(filename string, jobs chan<- Server) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening servers file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read the start of the JSON array
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("error reading json start: %v", err)
	}

	// Read each individual server object
	for decoder.More() {
		var server Server
		if err := decoder.Decode(&server); err != nil {
			return fmt.Errorf("error decoding json objects: %v", err)
		}
		jobs <- server
	}

	// Read the end of the JSON array
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("error reading json end: %v", err)
	}

	return nil
}

// LoadServersFull loads the entire server list from a small JSON file
func LoadServersFull(filename string) ([]Server, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening servers file: %v", err)
	}
	defer file.Close()

	var serversList []Server
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&serversList); err != nil {
		return nil, fmt.Errorf("error while decoding servers file: %v", err)
	}

	return serversList, nil
}
