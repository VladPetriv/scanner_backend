package handler

import (
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type HomePageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	Messages       []model.FullMessage
	ChannelsLength int
	MessagesLength int
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	data := HomePageData{
		Title: "Telegram Overflow",
		Type:  "messages",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	messages, err := h.service.Message.GetFullMessages()
	if err != nil {
		h.log.Error(err)
	}

	data.Channels = channels
	data.Messages = messages
	data.ChannelsLength = len(channels)
	data.MessagesLength = len(messages)

	h.tmpTree["messages"] = template.Must(
		template.ParseFiles("templates/messages.html", "templates/navbar.html", "templates/header.html", "templates/channels.html", "templates/channel.html", "templates/base.html"),
	)
	h.tmpTree["messages"].ExecuteTemplate(w, "base", data)
}
