package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/gorilla/mux"
)

type UserPageData struct {
	Type           string
	Title          string
	User           model.User
	Channels       []model.Channel
	ChannelsLength int
	Messages       []model.FullMessage
	MessagesLength int
}

func (h *Handler) userPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	data := UserPageData{
		Type:  "user",
		Title: "Telegram User",
	}

	ID, _ := strconv.Atoi(userID)

	user, err := h.service.User.GetUserByID(ID)
	if err != nil {
		h.log.Error(err)
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	messages, err := h.service.Message.GetFullMessagesByUserID(user.ID)
	if err != nil {
		h.log.Error(err)
	}

	data.User = *user
	data.Channels = channels
	data.ChannelsLength = len(channels)
	data.Messages = messages
	data.MessagesLength = len(messages)

	h.tmpTree["user"] = template.Must(
		template.ParseFiles("templates/user.html", "templates/navbar.html", "templates/header.html", "templates/messages.html", "templates/channels.html", "templates/channel.html", "templates/base.html"),
	)
	h.tmpTree["user"].ExecuteTemplate(w, "base", data)
}
