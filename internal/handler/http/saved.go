package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type savedPageData struct {
	DefaultPageData PageData
	Messages        []model.FullMessage
	MessagesLength  int
}

func (h Handler) loadSavedMessagesPage(w http.ResponseWriter, r *http.Request) {
	data := savedPageData{
		DefaultPageData: PageData{
			Type:         "saved",
			Title:        "Saved user messages",
			WebUserEmail: "",
			WebUserID:    0,
		},
	}

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert user id to int")

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

	pageData, err := h.service.Saved.ProcessSavedMessages(userID)
	if err != nil {
		h.log.Error().Err(err).Msg("get data for saved page")
	}
	if pageData != nil {
		data.Messages = pageData.SavedMessages
		data.MessagesLength = pageData.SavedMessagesCount
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load saved messages page")
	}
}

func (h Handler) createSavedMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert user id to int")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	messageID, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert message id to int")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	err = h.service.Saved.CreateSavedMessage(&model.Saved{WebUserID: userID, MessageID: messageID})
	if err != nil {
		h.log.Error().Err(err).Msg("create saved message")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%v", user.ID), http.StatusFound)
}

func (h Handler) deleteSavedMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["saved_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("convert saved message id to int")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")

		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	err = h.service.Saved.DeleteSavedMessage(messageID)
	if err != nil {
		h.log.Error().Err(err).Msg("delete saved message")
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%d", user.ID), http.StatusMovedPermanently)
}
