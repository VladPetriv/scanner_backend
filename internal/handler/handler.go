package handler

import "github.com/gorilla/mux"

type Handler struct {
}

func NewHanlder() *Handler {
	return &Handler{}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", h.hello).Methods("GET")
	return router
}
