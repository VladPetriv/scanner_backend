package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type SavedPageData struct {
	DefaultPageData PageData
	Messages        []model.FullMessage
	MessagesLength  int
}

func (h Handler) savedPage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert user id to int")
	}

	messages := make([]model.FullMessage, 0)

	savedMessages, err := h.service.Saved.GetSavedMessages(userID)
	if err != nil {
		h.log.Error().Err(err).Msg("get saved messages")
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for navbar")
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	for _, msg := range savedMessages {
		fullMessage, err := h.service.Message.GetFullMessageByMessageID(msg.MessageID)
		if err != nil {
			h.log.Error().Err(err).Msg("get full message by message id")
		}

		fullMessage.SavedID = msg.ID

		messages = append(messages, *fullMessage)
	}

	data := SavedPageData{
		DefaultPageData: PageData{
			Type:           "saved",
			Title:          "Saved user messages",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   "",
			WebUserID:      0,
		},
		Messages:       messages,
		MessagesLength: len(messages),
	}
	if user != nil {
		data.DefaultPageData.WebUserEmail = user.Email
		data.DefaultPageData.WebUserID = user.ID
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("execute save page")
	}
}

func (h Handler) createSavedMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert user id to int")
	}
	messageID, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert message id to int")
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	err = h.service.Saved.CreateSavedMessage(&model.Saved{WebUserID: userID, MessageID: messageID})
	if err != nil {
		h.log.Error().Err(err).Msg("create saved message")

		http.Redirect(w, r, "/home", http.StatusConflict)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%d", user.ID), http.StatusFound)
}

func (h Handler) deleteSavedMessage(w http.ResponseWriter, r *http.Request) {
	savedID, err := strconv.Atoi(mux.Vars(r)["saved_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert saved message id to int")
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	err = h.service.Saved.DeleteSavedMessage(savedID)
	if err != nil {
		h.log.Error().Err(err).Msg("delete saved message")
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%d", user.ID), http.StatusMovedPermanently)
}
