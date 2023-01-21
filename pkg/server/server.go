package server

import (
	"fmt"
	"net/http"
	"time"
)

// operationTimeout is used to define which time server will have for read/write timeouts.
const operationTimeout = 10

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  operationTimeout * time.Second,
		WriteTimeout: operationTimeout * time.Second,
	}

	err := s.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}
