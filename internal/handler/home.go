package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladPetriv/go-pagination-bootstrap"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type HomePageData struct {
	DefaultPageData PageData
	Messages        []model.FullMessage
	MessagesLength  int
	Pager           *pagination.Pagination
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	messagesLength, err := h.service.Message.GetMessagesLength()
	if err != nil {
		h.log.Error(err)
	}

	messages, err := h.service.Message.GetFullMessages(page)
	if err != nil {
		h.log.Error(err)
	}

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	messages = checkMessagesStatus(messages, h.service)

	webUserID, webUserEmail := util.ProcessWebUserData(user)

	data := HomePageData{
		DefaultPageData: PageData{
			Title:          "Telegram Overflow",
			Type:           "messages",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   webUserEmail,
			WebUserID:      webUserID,
		},
		Messages:       messages,
		MessagesLength: messagesLength,
		Pager:          pagination.New(messagesLength, 10, page, "/home/?page=0"),
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
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
