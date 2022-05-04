package handler

import (
	"html/template"
	"net/http"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/gorilla/mux"
)

type ChannelPageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	ChannelsLength int
}
type SingleChannelPageData struct {
	Type           string
	Title          string
	Channel        model.Channel
	Channels       []model.Channel
	ChannelsLength int
	Messages       []model.FullMessage
	MessagesLength int
}

func (h *Handler) channelsPage(w http.ResponseWriter, r *http.Request) {
	data := ChannelPageData{
		Title: "Telegram channels",
		Type:  "channels",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	data.Channels = channels
	data.ChannelsLength = len(channels)

	h.tmpTree["channels"] = template.Must(
		template.ParseFiles(
			"templates/channels.html", "templates/navbar.html", "templates/header.html", "templates/message.html",
			"templates/messages.html", "templates/channel.html", "templates/user.html", "templates/base.html",
		),
	)
	h.tmpTree["channels"].ExecuteTemplate(w, "base", data)
}

func (h *Handler) channelPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelName := vars["channel_name"]

	data := SingleChannelPageData{
		Type:  "channel",
		Title: "Telegram channel",
	}
	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	channel, err := h.service.Channel.GetChannelByName(channelName)
	if err != nil {
		h.log.Error(err)
	}

	fullMessages, err := h.service.Message.GetFullMessagesByChannelID(channel.ID)
	if err != nil {
		h.log.Error(err)
	}

	data.Channel = *channel
	data.Channels = channels
	data.Messages = fullMessages
	data.ChannelsLength = len(channels)
	data.MessagesLength = len(fullMessages)

	h.tmpTree["singleChannel"] = template.Must(
		template.ParseFiles(
			"templates/channel.html", "templates/navbar.html", "templates/header.html", "templates/message.html",
			"templates/messages.html", "templates/channels.html", "templates/user.html", "templates/base.html",
		),
	)
	h.tmpTree["singleChannel"].ExecuteTemplate(w, "base", data)
}
