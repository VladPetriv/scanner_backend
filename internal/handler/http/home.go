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
	data := homePageData{
		DefaultPageData: PageData{
			Title:        "Telegram Overflow",
			Type:         "messages",
			WebUserEmail: "",
			WebUserID:    0,
		},
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error().Err(err).Msg("convert page to int")
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for nav bar")
	}
	if navBarChannels != nil {
		data.DefaultPageData.Channels = GetRightChannelsCountForNavBar(navBarChannels)
		data.DefaultPageData.ChannelsLength = len(navBarChannels)
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}
	if user != nil {
		data.DefaultPageData.WebUserEmail = user.Email
		data.DefaultPageData.WebUserID = user.ID
	}

	pageData, err := h.service.Message.ProcessHomePage(page)
	if err != nil {
		h.log.Error().Err(err).Msg("process home page")
	}
	if pageData != nil {
		pageData.Messages = updateMessagesStatuses(pageData.Messages, h.service)

		data.Messages = pageData.Messages
		data.MessagesLength = pageData.MessagesCount
		data.Pager = pagination.New(pageData.MessagesCount, messagesPerPage, page, "/home/?page=0")
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
