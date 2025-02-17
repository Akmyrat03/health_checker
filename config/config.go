package config

import (
	"bufio"
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

func LoadServers(filename string) ([]string, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("error getting file info %v", err)
	}

	if info.Size() > limit {
		return LoadServersLineByLine(filename)
	} else {
		return LoadServersFull(filename)
	}
}

func LoadServersLineByLine(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening servers file: %v", err)
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)

	arrayStarted := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if !arrayStarted && line == "[" {
			arrayStarted = true
			continue
		}

		if arrayStarted && line == "]" {
			break
		}

		var server struct {
			Server string `json:"server"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &server); err != nil {
			return nil, fmt.Errorf("error during server JSON: %v", err)
		}
		servers = append(servers, server.Server)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading servers file: %v", err)
	}

	return servers, nil
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
