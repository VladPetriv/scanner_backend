package main

import (
	"log"

	"github.com/VladPetriv/scanner_backend/config"
	"github.com/VladPetriv/scanner_backend/internal/handler"
	"github.com/VladPetriv/scanner_backend/internal/server"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("error while getting conifg: %v", err)
	}

	srv := new(server.Server)

	handler := handler.NewHanlder()

	log.Printf("starting server at bind addr %s", cfg.BindAddr)
	if err := srv.Run(cfg.BindAddr, handler.InitRouter()); err != nil {
		log.Fatal("error while starting server")
	}
}
