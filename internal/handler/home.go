package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/AndyEverLie/go-pagination-bootstrap"
	"github.com/VladPetriv/scanner_backend/internal/model"
)

type HomePageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	Messages       []model.FullMessage
	ChannelsLength int
	MessagesLength int
	Pager          *pagination.Pagination
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	iPage, _ := strconv.Atoi(page)

	data := HomePageData{
		Title: "Telegram Overflow",
		Type:  "messages",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	messagesLength, err := h.service.Message.GetMessages()
	if err != nil {
		h.log.Error(err)
	}

	messages, err := h.service.Message.GetFullMessages(iPage)
	if err != nil {
		h.log.Error(err)
	}

	pager := pagination.New(len(messagesLength), 10, iPage, "/home?page=0")

	data.Channels = channels[:10]
	data.Messages = messages
	data.ChannelsLength = len(channels)
	data.MessagesLength = len(messagesLength)
	data.Pager = pager

	h.tmpTree["messages"] = template.Must(
		template.ParseFiles(
			"templates/message/messages.html", "templates/partials/navbar.html", "templates/partials/header.html", "templates/message/message.html",
			"templates/channel/channels.html", "templates/channel/channel.html", "templates/user/user.html", "templates/base.html",
		),
	)
	h.tmpTree["messages"].ExecuteTemplate(w, "base", data)
}
