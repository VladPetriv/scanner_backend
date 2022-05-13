package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type PageData struct {
	Title   string
	Message string
}

func (h *Handler) registrationPage(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Registration",
	}

	h.tmpTree["register"] = template.Must(template.ParseFiles("templates/auth/register.html"))
	err := h.tmpTree["register"].Execute(w, data)
	if err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) loginPage(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "login",
	}

	h.tmpTree["login"] = template.Must(template.ParseFiles("templates/auth/login.html"))
	err := h.tmpTree["login"].Execute(w, data)
	if err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) registration(w http.ResponseWriter, r *http.Request) {
	u := h.getUserFromForm(r)

	candidate, err := h.service.WebUser.GetWebUserByEmail(u.Email)
	if candidate != nil && err == nil {
		h.tmpTree["register"].Execute(
			w,
			PageData{Title: "Registration", Message: fmt.Sprintf("User with email %s is exist", u.Email)},
		)
	}
	hashedPassword, _ := util.HashPassword(u.Password)

	u.Password = hashedPassword

	err = h.service.WebUser.CreateWebUser(u)
	if err != nil {
		h.tmpTree["register"].Execute(
			w,
			PageData{Title: "Registration", Message: "error while creating user.Please try again later!"},
		)
		h.log.Error(err)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	u := h.getUserFromForm(r)

	candidate, err := h.service.WebUser.GetWebUserByEmail(u.Email)
	if err != nil {
		h.log.Error(err)
	}

	if candidate == nil {
		h.tmpTree["login"].Execute(
			w,
			PageData{Title: "Login", Message: fmt.Sprintf("User with email %s not found.Maybe you want to signup?", u.Email)},
		)
	}

	if util.ComparePassword(u.Password, candidate.Password) {
		h.writeToSessionStore(w, r, u.Email)

		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	h.tmpTree["login"].Execute(
		w,
		PageData{Title: "Login", Message: "User password is incorrect!"},
	)
}
