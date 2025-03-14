package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	cfg  Config
	once sync.Once
)

type Config struct {
	SMTP     SMTP       `env:"SMTP_CONFIG"`
	App      App        `env:"APP"`
	Postgres PostgresDB `env:"POSTGRES"`
	Cors     Cors       `env:"CORS"`
	JWT      JWT        `env:"JWT"`
}

type PostgresDB struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
	SslMode  string `env:"POSTGRES_SSLMODE" env-required:"true"`
}

type SMTP struct {
	SMTPServer    string `env:"SMTP_SERVER" env-required:"true"`
	SMTPPort      int    `env:"SMTP_PORT" env-required:"true"`
	SMTPEmail     string `env:"SMTP_EMAIL" env-required:"true"`
	SMTPPass      string `env:"SMTP_PASS" env-required:"true"`
	SubjectPrefix string `env:"SMTP_SUBJECT_PREFIX" env-required:"true"`
}

type App struct {
	Host     string `env:"APP_HOST" env-required:"true"`
	Port     string `env:"APP_PORT" env-required:"true"`
	EndPoint string `env:"APP_ENDPOINT" env-required:"true"`
}

type Cors struct {
	Origins     string `env:"CORS_ORIGINS" env-required:"true"`
	Credentials bool   `env:"CORS_CREDENTIALS" env-required:"true"`
}

type JWT struct {
	JwtSecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
}

func LoadConfig() *Config {
	once.Do(func() {
		err := cleanenv.ReadConfig(".env", &cfg)
		if err != nil {
			log.Fatalf("Error loading configuration: %v", err)
		}

		log.Println("Configuration loaded successfully")
	})

	return &cfg
}
