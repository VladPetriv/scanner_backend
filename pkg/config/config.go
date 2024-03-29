package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PgUser         string
	PgPassword     string
	PgDB           string
	PgHost         string
	MigrationsPath string
	Port           string
	DatabaseURL    string
	LogLevel       string
	LogFilename    string
	KafkaAddr      string
	CookieSecret   string
}

func Get() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("get .env file error: %w", err)
	}

	return &Config{
		PgUser:         os.Getenv("POSTGRES_USER"),
		PgPassword:     os.Getenv("POSTGRES_PASSWORD"),
		PgDB:           os.Getenv("POSTGRES_DB"),
		PgHost:         os.Getenv("POSTGRES_HOST"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		Port:           os.Getenv("PORT"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		LogLevel:       os.Getenv("LOG_LEVEL"),
		LogFilename:    os.Getenv("LOG_FILENAME"),
		KafkaAddr:      os.Getenv("KAFKA_ADDR"),
		CookieSecret:   os.Getenv("COOKIE_SECRET"),
	}, nil
}
