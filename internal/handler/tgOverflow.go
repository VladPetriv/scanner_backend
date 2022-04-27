package handler

import (
	"net/http"
)

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
