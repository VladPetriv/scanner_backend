package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/VladPetriv/go-pagination-bootstrap"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type HomePageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	Messages       []model.FullMessage
	ChannelsLength int
	MessagesLength int
	Pager          *pagination.Pagination
	UserEmail      interface{}
	WebUserID      int
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

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	messages = checkMessagesStatus(messages, h.service)

	pager := pagination.New(len(messagesLength), 10, iPage, "/home?page=0")

	data.Channels = util.ProcessChannels(channels)
	data.Messages = messages
	data.ChannelsLength = len(channels)
	data.MessagesLength = len(messagesLength)
	data.Pager = pager
	data.WebUserID, data.UserEmail = util.ProcessWebUserData(user)

	h.tmpTree["messages"] = template.Must(
		template.ParseFiles(
			"templates/message/messages.html", "templates/partials/navbar.html", "templates/partials/header.html", "templates/message/message.html",
			"templates/channel/channels.html", "templates/channel/channel.html", "templates/user/saved.html", "templates/user/user.html",
			"templates/base.html",
		),
	)
	err = h.tmpTree["messages"].ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error(err)
	}
}

func checkMessagesStatus(messages []model.FullMessage, manager *service.Manager) []model.FullMessage {
	result := make([]model.FullMessage, 0)

	for _, message := range messages {
		saved, err := manager.Saved.GetSavedMessageByMessageID(message.ID)
		if err == nil && saved != nil {
			message.Status = true
			result = append(result, message)
			continue
		}

		message.Status = false
		result = append(result, message)
	}

	return result
}
