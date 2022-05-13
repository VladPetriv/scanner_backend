package handler

import (
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Handler struct {
	store   *sessions.CookieStore
	service *service.Manager
	log     *logger.Logger
	tmpTree map[string]*template.Template
}

func NewHandler(serviceManager *service.Manager, log *logger.Logger) *Handler {
	return &Handler{
		store:   sessions.NewCookieStore([]byte("secret")),
		service: serviceManager,
		log:     log,
		tmpTree: make(map[string]*template.Template),
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/home", http.StatusFound)).Methods("GET")
	router.HandleFunc("/home", h.homePage).Methods("GET")
	router.HandleFunc("/channel", h.channelsPage).Methods("GET")
	router.HandleFunc("/channel/{channel_name}", h.channelPage).Methods("GET")
	router.HandleFunc("/user/{user_id}", h.userPage).Methods("GET")
	router.HandleFunc("/message/{message_id}", h.messagePage).Methods("GET")
	router.HandleFunc("/registration", h.registrationPage).Methods("GET")
	router.HandleFunc("/login", h.loginPage).Methods("GET")
	router.HandleFunc("/registration", h.registration).Methods("POST")
	router.HandleFunc("/login", h.login).Methods("POST")

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

func (h *Handler) checkUserStatus(r *http.Request) interface{} {
	sessions, err := h.store.Get(r, "session")
	if err != nil {
		h.log.Error("Session error: ", err)
	}

	email, ok := sessions.Values["userEmail"]
	if ok {
		return email
	}

	return ""
}

func (h *Handler) writeToSessionStore(w http.ResponseWriter, r *http.Request, value interface{}) {
	session, _ := h.store.Get(r, "session")
	session.Values["userEmail"] = value
	session.Save(r, w)
}

func (h *Handler) getUserFromForm(r *http.Request) *model.WebUser {
	r.ParseForm()

	user := &model.WebUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	return user
}
