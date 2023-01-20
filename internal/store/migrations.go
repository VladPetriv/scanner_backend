package store

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	// postgres and file module are required for running migration.
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/VladPetriv/scanner_backend/pkg/config"
)

func runMigrations(cfg *config.Config) error {
	var connectionString string

	if cfg.MigrationsPath == "" {
		return nil
	}

	if cfg.DatabaseURL == "" {
		connectionString = fmt.Sprintf(
			"postgresql://%s/%s?user=%s&password=%s&sslmode=disable",
			cfg.PgHost, cfg.PgDB, cfg.PgUser, cfg.PgDB,
		)
	} else {
		connectionString = cfg.DatabaseURL
	}

	m, err := migrate.New(
		cfg.MigrationsPath,
		connectionString,
	)
	if err != nil {
		return fmt.Errorf("create migrations error: %w", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("up migrations error: %w", err)
		}
	}

	return nil
}
