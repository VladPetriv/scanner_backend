package main

import (
	_ "github.com/lib/pq"

	"github.com/VladPetriv/scanner_backend/internal/handler"
	"github.com/VladPetriv/scanner_backend/internal/server"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

func main() {
	log := logger.Get()

	cfg, err := config.Get()
	if err != nil {
		log.Errorf("error while getting config: %v", err)
	}

	store, err := store.New(cfg, log)
	if err != nil {
		log.Errorf("error while creating store: %v", err)
	}

	serviceManger, err := service.NewManager(store)
	if err != nil {
		log.Errorf("error while creating service manager: %v", err)
	}

	srv := new(server.Server)

	handler := handler.NewHandler(serviceManger, log)

	log.Infof("starting server at port %s", cfg.Port)

	if err := srv.Run(cfg.Port, handler.InitRouter()); err != nil {
		log.Error("error while starting server: ", err)
	}
}
