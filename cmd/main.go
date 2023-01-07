package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/VladPetriv/scanner_backend/internal/handler"
	"github.com/VladPetriv/scanner_backend/internal/server"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/kafka"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("get config error: %v", err)
	}

	log := logger.Get(cfg)

	store, err := store.New(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("create store")
	}

	serviceManger, err := service.NewManager(store)
	if err != nil {
		log.Fatal().Err(err).Msg("create service manager")
	}

	go kafka.SaveChannelsFromQueueToDB(serviceManger, cfg, log)
	go kafka.SaveDataFromQueueToDB(serviceManger, cfg, log)

	srv := new(server.Server)

	handler := handler.NewHandler(serviceManger, log)

	log.Info().Msgf("starting server at port: %s", cfg.Port)

	if err = srv.Run(cfg.Port, handler.InitRouter()); err != nil {
		log.Fatal().Err(err).Msg("start server")
	}
}
