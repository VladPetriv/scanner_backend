package handler

import (
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type HomePageData struct {
	Title          string
	Channels       []model.Channel
	Messages       []model.Message
	ChannelsLength int
	MessagesLength int
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	data := HomePageData{
		Title: "Telegram Overflow",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	messages, err := h.service.Message.GetMessages()
	if err != nil {
		h.log.Error(err)
	}

	data.Channels = channels
	data.Messages = messages
	data.ChannelsLength = len(channels)
	data.MessagesLength = len(messages)

	tmpTree := make(map[string]*template.Template)

	tmpTree["channels.html"] = template.Must(template.ParseFiles("templates/channels.html", "templates/messages.html", "templates/navbar.html", "templates/base.html"))

	tmpTree["channels.html"].ExecuteTemplate(w, "base", data)
}
