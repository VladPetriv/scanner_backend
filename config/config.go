package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PgUser         string
	PgPassword     string
	PgDb           string
	PgHost         string
	MigrationsPath string
	BindAddr       string
}

func Get() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error while loading .env file: %s", err)
	}

	return &Config{
		PgUser:         os.Getenv("POSTGRES_USER"),
		PgPassword:     os.Getenv("POSTGRES_PASSWORD"),
		PgDb:           os.Getenv("POSTGRES_DB"),
		PgHost:         os.Getenv("POSTGRES_HOST"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		BindAddr:       os.Getenv("BIND_ADDR"),
	}, nil
}
