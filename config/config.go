package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Servers       []string `json:"servers"`
	CheckInterval int      `json:"check_interval"`
	LogFile       string   `json:"log_file"`
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
