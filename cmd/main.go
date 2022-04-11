package main

import (
	"github.com/VladPetriv/scanner_backend/config"
	"github.com/VladPetriv/scanner_backend/internal/handler"
	"github.com/VladPetriv/scanner_backend/internal/server"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/logger"
)

func main() {
	log := logger.Get()

	cfg, err := config.Get()
	if err != nil {
		log.Errorf("error while getting conifg: %v", err)
	}

	store, err := store.New(*cfg, log)
	if err != nil {
		log.Errorf("error while creating store: %v", err)
	}

	serviceManger, err := service.NewManager(store)
	if err != nil {
		log.Errorf("error while creating service manager: %v", err)
	}

	srv := new(server.Server)

	handler := handler.NewHandler(serviceManger)

	log.Infof("starting server at bind addr %s", cfg.BindAddr)
	if err := srv.Run(cfg.BindAddr, handler.InitRouter()); err != nil {
		log.Error("error while starting server")
	}
}
