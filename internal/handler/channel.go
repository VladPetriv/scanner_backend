package handler

import (
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type ChannelPageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	ChannelsLength int
}

func (h *Handler) channelPage(w http.ResponseWriter, r *http.Request) {
	data := ChannelPageData{
		Title: "Telegram channels",
		Type:  "channel",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	data.Channels = channels
	data.ChannelsLength = len(channels)

	h.tmpTree["channels"] = template.Must(template.ParseFiles("templates/channel.html", "templates/navChannels.html", "templates/navbar.html", "templates/home.html", "templates/base.html"))
	h.tmpTree["channels"].ExecuteTemplate(w, "base", data)
}
