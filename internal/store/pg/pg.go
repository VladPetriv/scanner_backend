package pg

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func Init(cfg *config.Config) (*DB, error) {
	var connectionString string

	if cfg.DatabaseURL == "" {
		connectionString = fmt.Sprintf(
			"user=%s password=%s host=%s dbname=%s sslmode=disable",
			cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgDB,
		)
	} else {
		connectionString = cfg.DatabaseURL
	}

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error while create connection to db: %w", err)
	}
	_, err = db.Exec("SELECT 1;")
	if err != nil {
		return nil, fmt.Errorf("error while send request to db: %w", err)
	}

	return &DB{db}, nil
}
