package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
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
	data.Channels = util.ProcessChannels(channels)
	data.ChannelsLength = len(channels)
	data.Messages = messages
	data.MessagesLength = len(messages)

	h.tmpTree["user"] = template.Must(
		template.ParseFiles(
			"templates/user/user.html", "templates/partials/navbar.html", "templates/partials/header.html", "templates/message/message.html",
			"templates/message/messages.html", "templates/channel/channels.html", "templates/channel/channel.html", "templates/base.html",
		),
	)
	err = h.tmpTree["user"].ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error(err)
	}
}
