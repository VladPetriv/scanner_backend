package handler

import (
	"net/http"
	"strconv"

	"github.com/VladPetriv/go-pagination-bootstrap"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
)

const messagesPerPage = 10

type homePageData struct {
	DefaultPageData PageData
	Messages        []model.FullMessage
	MessagesLength  int
	Pager           *pagination.Pagination
}

func (h Handler) loadHomePage(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error().Err(err).Msg("convert page to int")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for nav bar")
	}

	messagesCount, err := h.service.Message.GetMessagesCount()
	if err != nil {
		h.log.Error().Err(err).Msg("get messages count")
	}

	messages, err := h.service.Message.GetFullMessagesByPage(page)
	if err != nil {
		h.log.Error().Err(err).Msg("get full messages by page")
	}

	messages = updateMessagesStatuses(messages, h.service)

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	data := homePageData{
		DefaultPageData: PageData{
			Title:          "Telegram Overflow",
			Type:           "messages",
			Channels:       GetRightChannelsCountForNavBar(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   "",
			WebUserID:      0,
		},
		Messages:       messages,
		MessagesLength: messagesCount,
		Pager:          pagination.New(messagesCount, messagesPerPage, page, "/home/?page=0"),
	}
	if user != nil {
		data.DefaultPageData.WebUserEmail = user.Email
		data.DefaultPageData.WebUserID = user.ID
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load home page")
	}
}

func updateMessagesStatuses(messages []model.FullMessage, manager *service.Manager) []model.FullMessage {
	var result []model.FullMessage

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
