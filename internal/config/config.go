package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	SMTP     SMTP       `json:"smtp_config"`
	App      App        `json:"app"`
	Postgres PostgresDB `json:"postgres"`
	Cors     Cors       `json:"cors"`
	JWT      JWT        `json:"jwt"`
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
	SMTPServer    string `json:"smtp_server"`
	SMTPPort      int    `json:"smtp_port"`
	SMTPEmail     string `json:"smtp_email"`
	SMTPPass      string `json:"smtp_pass"`
	SubjectPrefix string `json:"subject_prefix"`
}

type App struct {
	Host     string `json:"app_host"`
	Port     string `json:"app_port"`
	EndPoint string `json:"app_endpoint"`
}

type Cors struct {
	Origins     string `json:"cors_origins"`
	Credentials bool   `json:"cors_credentials"`
}

type JWT struct {
	JwtSecretKey string `json:"jwt_secret_key"`
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
