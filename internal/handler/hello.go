package handler

import (
	"net/http"
)

func (h *Handler) hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
