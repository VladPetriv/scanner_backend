package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type MessagePageData struct {
	DefaultPageData PageData
	Message         model.FullMessage
}

func (h Handler) loadMessagePage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert message_id to int")
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for navbar")
	}

	message, err := h.service.Message.GetFullMessageByMessageID(id)
	if err != nil {
		h.log.Error().Err(err).Msg("get full messages by message id")
	}

	replies, err := h.service.Reply.GetFullRepliesByMessageID(message.ID)
	if err != nil {
		h.log.Error().Err(err).Msg("get full replies by message id")
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	message.Replies = replies

	data := MessagePageData{
		DefaultPageData: PageData{
			Type:           "message",
			Title:          "Telegram message",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   "",
			WebUserID:      0,
		},
		Message: *message,
	}
	if user != nil {
		data.DefaultPageData.WebUserEmail = user.Email
		data.DefaultPageData.WebUserID = user.ID
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load message page")
	}
}
