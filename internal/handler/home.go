package handler

import (
	"net/http"
)

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	h.ProcessTemplateData(w, "templates/base.html", "Home page")
}
