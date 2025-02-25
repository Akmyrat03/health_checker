package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Basic    Basic      `json:"basic"`
	SMTP     SMTP       `json:"smtp_config"`
	Servers  []Server   `json:"servers"`
	App      App        `json:"app"`
	Postgres PostgresDB `json:"postgres"`
}

type Basic struct {
	Interval int `json:"interval"`
	Timeout  int `json:"timeout"`
}

type PostgresDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SslMode  string `json:"sslmode"`
}

type SMTP struct {
	SMTPServer    string   `json:"smtp_server"`
	SMTPPort      int      `json:"smtp_port"`
	SMTPEmail     string   `json:"smtp_email"`
	SMTPPass      string   `json:"smtp_pass"`
	SubjectPrefix string   `json:"subject_prefix"`
	Receivers     []string `json:"receivers"`
}

type Server struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type App struct {
	Host     string `json:"app_host"`
	Port     string `json:"app_port"`
	EndPoint string `json:"app_endpoint"`
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
