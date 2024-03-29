package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/internal/service"
)

type authPageData struct {
	Title   string
	Message string
}

func (h Handler) loadRegistrationPage(w http.ResponseWriter, r *http.Request) {
	data := authPageData{
		Title: "Registration",
	}

	h.tmpTree["register"] = template.Must(template.ParseFiles("templates/auth/register.html"))
	err := h.tmpTree["register"].Execute(w, data)
	if err != nil {
		h.log.Error().Err(err).Msg("load registration page")
	}
}

func (h Handler) loadLoginPage(w http.ResponseWriter, r *http.Request) {
	data := authPageData{
		Title: "login",
	}

	h.tmpTree["login"] = template.Must(template.ParseFiles("templates/auth/login.html"))
	err := h.tmpTree["login"].Execute(w, data)
	if err != nil {
		h.log.Error().Err(err).Msg("load login page")
	}
}

func (h Handler) registration(w http.ResponseWriter, r *http.Request) {
	data := authPageData{
		Title: "Registration",
	}

	user := h.getUserFromForm(r)

	err := h.service.Auth.Register(user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWebUserIsExist):
			data.Message = fmt.Sprintf("User with email %s is exist!", user.Email)
		default:
			data.Message = "Failed to register new user!"
		}
		err = h.tmpTree["register"].Execute(w, data)
		if err != nil {
			h.log.Error().Err(err).Msg("execute register template")
		}
	}

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	data := authPageData{
		Title: "Login",
	}

	user := h.getUserFromForm(r)

	email, err := h.service.Auth.Login(user.Email, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWebUserNotFound):
			data.Message = fmt.Sprintf("User with email %s not found!", user.Email)
		case errors.Is(err, service.ErrIncorrectPassword):
			data.Message = "User password is incorrect!"
		default:
			data.Message = "Failed to login!"

			h.log.Error().Err(err).Msg("login user")
		}
	}

	if email != "" {
		h.addUserToSession(w, r, email)

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)

		return
	}

	err = h.tmpTree["login"].Execute(w, data)
	if err != nil {
		h.log.Error().Err(err).Msg("execute login template")
	}
}

func (h Handler) logout(w http.ResponseWriter, r *http.Request) {
	h.removeUserFromSession(w, r)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
