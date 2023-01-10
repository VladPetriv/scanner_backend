package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/pkg/password"
)

type AuthPageData struct {
	Title   string
	Message string
}

func (h *Handler) registrationPage(w http.ResponseWriter, r *http.Request) {
	data := AuthPageData{
		Title: "Registration",
	}

	h.tmpTree["register"] = template.Must(template.ParseFiles("templates/auth/register.html"))
	err := h.tmpTree["register"].Execute(w, data)
	if err != nil {
		h.log.Error().Err(err).Msg("execute registration templates")
	}
}

func (h *Handler) loginPage(w http.ResponseWriter, r *http.Request) {
	data := AuthPageData{
		Title: "login",
	}

	h.tmpTree["login"] = template.Must(template.ParseFiles("templates/auth/login.html"))
	err := h.tmpTree["login"].Execute(w, data)
	if err != nil {
		h.log.Error().Err(err).Msg("execute login templates")
	}
}

func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	u := h.getUserFromForm(r)

	candidate, err := h.service.WebUser.GetWebUserByEmail(u.Email)
	if candidate != nil && err == nil {
		h.tmpTree["register"].Execute(
			w,
			AuthPageData{Title: "Registration", Message: fmt.Sprintf("User with email %s is exist", u.Email)},
		)
	}

	hashedPassword, err := password.HashPassword(u.Password)
	if err != nil {
		h.log.Error().Err(err).Msg("hash password")
	}

	u.Password = hashedPassword

	err = h.service.WebUser.CreateWebUser(u)
	if err != nil {
		h.tmpTree["register"].Execute(
			w,
			AuthPageData{Title: "Registration", Message: "Error while creating user. Please try again later!"},
		)
		h.log.Error().Err(err).Msg("execute registration template")
	}

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	u := h.getUserFromForm(r)

	candidate, err := h.service.WebUser.GetWebUserByEmail(u.Email)
	if err != nil {
		h.log.Error().Err(err).Msg("get user by email")
	}

	if candidate == nil {
		err = h.tmpTree["login"].Execute(
			w,
			AuthPageData{Title: "Login", Message: fmt.Sprintf("User with email %s not found", u.Email)},
		)
		if err != nil {
			h.log.Error().Err(err).Msg("execute login template")
		}
	}

	if password.ComparePassword(u.Password, candidate.Password) {
		h.writeToSessionStore(w, r, u.Email)

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	err = h.tmpTree["login"].Execute(
		w,
		AuthPageData{Title: "Login", Message: "User password is incorrect!"},
	)
	if err != nil {
		h.log.Error().Err(err).Msg("execute login template")
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	h.removeFromSessionStore(w, r)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
