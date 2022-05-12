package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/AndyEverLie/go-pagination-bootstrap"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/gorilla/mux"
)

type ChannelPageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	ChannelsLength int
	MainChannels   []model.Channel
	Pager          *pagination.Pagination
}

type SingleChannelPageData struct {
	Type           string
	Title          string
	Channel        model.Channel
	Channels       []model.Channel
	ChannelsLength int
	Messages       []model.FullMessage
	MessagesLength int
	Pager          *pagination.Pagination
}

func (h *Handler) channelsPage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	iPage, _ := strconv.Atoi(page)

	data := ChannelPageData{
		Title: "Telegram channels",
		Type:  "channels",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	mainChannels, err := h.service.Channel.GetChannelsByPage(iPage)
	if err != nil {
		h.log.Error(err)
	}

	pager := pagination.New(len(channels), 10, iPage, "/channel?=0")

	data.Channels = channels[:10]
	data.ChannelsLength = len(channels)
	data.MainChannels = mainChannels
	data.Pager = pager

	h.tmpTree["channels"] = template.Must(
		template.ParseFiles(
			"templates/channel/channels.html", "templates/partials/navbar.html", "templates/partials/header.html", "templates/message/message.html",
			"templates/message/messages.html", "templates/channel/channel.html", "templates/user/user.html", "templates/base.html",
		),
	)
	h.tmpTree["channels"].ExecuteTemplate(w, "base", data)
}

func (h *Handler) channelPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelName := vars["channel_name"]

	page := r.URL.Query().Get("page")
	iPage, _ := strconv.Atoi(page)

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

	length, err := h.service.Message.GetFullMessagesByChannelID(channel.ID, 1024*1024, iPage)
	if err != nil {
		h.log.Error(err)
	}

	fullMessages, err := h.service.Message.GetFullMessagesByChannelID(channel.ID, 10, iPage)
	if err != nil {
		h.log.Error(err)
	}

	pager := pagination.New(len(length), 10, iPage, "/channel/ru_python?page=0")

	data.Channel = *channel
	data.Channels = channels[:10]
	data.Messages = fullMessages
	data.ChannelsLength = len(channels)
	data.MessagesLength = len(length)
	data.Pager = pager

	h.tmpTree["singleChannel"] = template.Must(
		template.ParseFiles(
			"templates/channel/channel.html", "templates/partials/navbar.html", "templates/partials/header.html", "templates/message/message.html",
			"templates/message/messages.html", "templates/channel/channels.html", "templates/user/user.html", "templates/base.html",
		),
	)
	h.tmpTree["singleChannel"].ExecuteTemplate(w, "base", data)
}
