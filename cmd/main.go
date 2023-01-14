package main

import (
	"log"

	_ "github.com/lib/pq"

	handler "github.com/VladPetriv/scanner_backend/internal/handler/http"
	"github.com/VladPetriv/scanner_backend/internal/handler/queue/kafka"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
	"github.com/VladPetriv/scanner_backend/pkg/server"
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

	queue := kafka.New(serviceManger, cfg, log)
	go queue.SaveChannelsData()
	go queue.SaveMessagesData()

	srv := new(server.Server)

	httpHandler := handler.NewHandler(serviceManger, log)

	log.Info().Msgf("starting server at port: %s", cfg.Port)

	if err = srv.Run(cfg.Port, httpHandler.InitRouter()); err != nil {
		log.Fatal().Err(err).Msgf("start server at port: %s", cfg.Port)
	}
}
