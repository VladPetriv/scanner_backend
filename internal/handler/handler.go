package handler

import (
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Manager
}

func NewHandler(serviceManager *service.Manager) *Handler {
	return &Handler{
		service: serviceManager,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("tgOrverflow/", h.mainPage).Methods("GET")
	return router
}
