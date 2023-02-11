package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
)

type messagePageData struct {
	DefaultPageData PageData
	Message         model.FullMessage
}

func (h Handler) loadMessagePage(w http.ResponseWriter, r *http.Request) {
	data := messagePageData{
		DefaultPageData: PageData{
			Type:         "message",
			Title:        "Telegram message",
			WebUserEmail: "",
			WebUserID:    0,
		},
	}

	messageID, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert message_id to int")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for navbar")
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

	pageData, err := h.service.Message.ProcessMessagePage(messageID)
	if err != nil {
		if !errors.Is(err, service.ErrMessageNotFound) && !errors.Is(err, service.ErrRepliesNotFound) {
			h.log.Error().Err(err).Msg("process data for message page")
		}
	}
	if pageData != nil {
		data.Message = *pageData.Message
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load message page")
	}
}
