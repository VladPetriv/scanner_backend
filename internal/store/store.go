package store

import (
	"fmt"
	"time"

	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type Store struct {
	pg     *pg.DB
	logger *logger.Logger

	Channel ChannelRepo
	Message MessageRepo
	Replie  ReplieRepo
	User    UserRepo
	WebUser WebUserRepo
	Saved   SavedRepo
}

func New(cfg *config.Config, log *logger.Logger) (*Store, error) {
	pgDB, err := pg.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgresql: %w", err)
	}

	if pgDB != nil {
		log.Info().Msg("Running migrations...")
		if err := runMigrations(cfg); err != nil {
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	var store Store
	store.logger = log

	if pgDB != nil {
		store.pg = pgDB
		store.Channel = pg.NewChannelRepo(pgDB)
		store.Message = pg.NewMessageRepo(pgDB)
		store.Replie = pg.NewReplieRepo(pgDB)
		store.User = pg.NewUserRepo(pgDB)
		store.WebUser = pg.NewWebUserRepo(pgDB)
		store.Saved = pg.NewSavedRepo(pgDB)

		go store.KeepAliveDB(cfg)
	}

	return &store, nil
}

func (s *Store) KeepAliveDB(cfg *config.Config) {
	var err error

	for {
		time.Sleep(time.Second * 5)

		lostConnection := false
		if s.pg == nil {
			lostConnection = true
		} else if _, err := s.pg.Exec("SELECT 1;"); err != nil {
			lostConnection = true
		}

		if !lostConnection {
			continue
		}

		s.logger.Debug().Msg("[store.KeepAliveDB] Lost db connection. Restoring...")

		s.pg, err = pg.Init(cfg)
		if err != nil {
			s.logger.Error().Err(err).Msg("init postgresql")

			continue
		}

		s.logger.Debug().Msg("[store.KeepAliveDB] DB reconnected")
	}
}
