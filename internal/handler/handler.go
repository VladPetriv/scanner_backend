package handler

import (
	"html/template"

	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/logger"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Manager
	log     *logger.Logger
	tmpTree map[string]*template.Template
}

func NewHandler(serviceManager *service.Manager, log *logger.Logger) *Handler {
	return &Handler{
		service: serviceManager,
		log:     log,
		tmpTree: make(map[string]*template.Template),
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/home", h.homePage).Methods("GET")
	router.HandleFunc("/channel", h.channelPage).Methods("GET")

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			h.log.Error(err)
		}

		met, err := route.GetMethods()
		if err != nil {
			h.log.Error(err)
		}

		h.log.Infof("Route - %s %s", tpl, met)

		return nil
	})

	return router
}
