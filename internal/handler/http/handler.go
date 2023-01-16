package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type Handler struct {
	store     *sessions.CookieStore
	service   *service.Manager
	log       *logger.Logger
	tmpTree   map[string]*template.Template
	templates *template.Template
}

type PageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	ChannelsLength int
	WebUserEmail   interface{}
	WebUserID      int
}

func NewHandler(serviceManager *service.Manager, log *logger.Logger) *Handler {
	return &Handler{
		store:   sessions.NewCookieStore([]byte("secret")),
		service: serviceManager,
		log:     log,
		tmpTree: make(map[string]*template.Template),
		templates: template.Must(
			template.ParseFiles(
				"templates/message/messages.html", "templates/partials/navbar.html",
				"templates/partials/header.html", "templates/message/message.html",
				"templates/channel/channels.html", "templates/channel/channel.html",
				"templates/user/saved.html", "templates/user/user.html",
				"templates/base.html",
			),
		),
	}
}

func (h Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()

	home := router.PathPrefix("/").Subrouter()
	home.Handle("/", http.RedirectHandler("/home", http.StatusMovedPermanently)).Methods("GET")
	home.HandleFunc("/home", h.homePage).Methods("GET")

	channel := router.PathPrefix("/channel").Subrouter()
	channel.HandleFunc("/", h.channelsPage).Methods("GET")
	channel.HandleFunc("/{channel_name}", h.channelPage).Methods("GET")

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{user_id}", h.userPage).Methods("GET")

	message := router.PathPrefix("/message").Subrouter()
	message.HandleFunc("/{message_id}", h.messagePage).Methods("GET")

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", h.login).Methods("POST")
	auth.HandleFunc("/registration", h.registration).Methods("POST")
	auth.HandleFunc("/logout", h.logout).Methods("POST")
	auth.HandleFunc("/login", h.loadLoginPage).Methods("GET")
	auth.HandleFunc("/registration", h.loadRegistrationPage).Methods("GET")

	saved := router.PathPrefix("/saved").Subrouter()
	saved.HandleFunc("/{user_id}", h.savedPage).Methods("GET")
	saved.HandleFunc("/delete/{saved_id}", h.deleteSavedMessage).Methods("POST")
	saved.HandleFunc("/create/{user_id}/{message_id}", h.createSavedMessage).Methods("POST")

	h.logAllRoutes(router)

	return router
}

func (h Handler) logAllRoutes(router *mux.Router) {
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			h.log.Error().Err(err).Msg("get template path")
		}

		met, _ := route.GetMethods()
		if len(met) == 0 {
			met = append(met, "SUBROUTER")
		}

		h.log.Info().Msgf("Route - %s %s", tpl, met)

		return nil
	})
	if err != nil {
		h.log.Error().Err(err).Msg("walk through routes")
	}
}

func (h *Handler) checkUserStatus(r *http.Request) interface{} {
	sessions, err := h.store.Get(r, "session")
	if err != nil {
		h.log.Error().Err(err).Msg("get user session")
	}

	email, ok := sessions.Values["userEmail"]
	if ok {
		return email
	}

	return ""
}

func (h *Handler) writeToSessionStore(w http.ResponseWriter, r *http.Request, value interface{}) {
	session, err := h.store.Get(r, "session")
	if err != nil {
		h.log.Error().Err(err).Msg("get user session")
	}

	session.Values["userEmail"] = value
	err = session.Save(r, w)
	if err != nil {
		h.log.Error().Err(err).Msg("save user session")
	}
}

func (h *Handler) removeFromSessionStore(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "session")
	if err != nil {
		h.log.Error().Err(err).Msg("get user session")
	}

	delete(session.Values, "userEmail")

	err = session.Save(r, w)
	if err != nil {
		h.log.Error().Err(err).Msg("save session")
	}
}

func (h *Handler) getUserFromForm(r *http.Request) *model.WebUser {
	err := r.ParseForm()

	if err != nil {
		h.log.Error().Err(err).Msg("get form")
	}

	user := &model.WebUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	return user
}
