package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	limit = 10 * 1024 * 1024
)

type Config struct {
	CheckInterval int    `json:"check_interval"`
	LogFile       string `json:"log_file"`
	Timeout       int    `json:"timeout"`
}

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

func LoadServers(filename string) (<-chan string, error) {
	serversChan := make(chan string)

	go func() {
		defer close(serversChan)

		info, err := os.Stat(filename)
		if err != nil {
			fmt.Println("error getting file info:", err)
			return
		}

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

func LoadServersLineByLine(filename string, jobs chan<- string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening servers file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("error reading json start: %v", err)
	}

	for decoder.More() {
		var server struct {
			Server string `json:"server"`
		}
		if err := decoder.Decode(&server); err != nil {
			return fmt.Errorf("error decoding json objects: %v", err)
		}
		jobs <- server.Server
	}

	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("error reading json end: %v", err)
	}

	return nil
}

func LoadServersFull(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening servers file: %v", err)
	}
	defer file.Close()

	var serversList []struct {
		Server string `json:"server"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&serversList); err != nil {
		return nil, fmt.Errorf("error while decoding servers file: %v", err)
	}

	var servers []string
	for _, s := range serversList {
		servers = append(servers, s.Server)
	}

	return servers, nil
}
