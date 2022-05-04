package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/gorilla/mux"
)

type MessagePageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	ChannelsLength int
	Message        model.FullMessage
}

func (h *Handler) messagePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["message_id"]

	data := MessagePageData{
		Type:  "message",
		Title: "Telegram message",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	ID, _ := strconv.Atoi(messageID)

	message, err := h.service.Message.GetFullMessageByMessageID(ID)
	if err != nil {
		h.log.Error(err)
	}

	replies, err := h.service.Replie.GetFullRepliesByMessageID(message.ID)
	if err != nil {
		h.log.Error(err)
	}

	message.Replies = replies

	data.Channels = channels
	data.ChannelsLength = len(channels)
	data.Message = *message

	h.tmpTree["channels"] = template.Must(
		template.ParseFiles(
			"templates/channels.html", "templates/navbar.html", "templates/header.html", "templates/message.html",
			"templates/messages.html", "templates/channel.html", "templates/user.html", "templates/base.html",
		),
	)
	h.tmpTree["channels"].ExecuteTemplate(w, "base", data)
}
