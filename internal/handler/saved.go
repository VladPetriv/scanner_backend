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

func (h *Handler) savedPage(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["user_id"])

	messages := make([]model.FullMessage, 0)

	savedMessages, err := h.service.Saved.GetSavedMessages(userID)
	if err != nil {
		h.log.Error(err)
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	for _, msg := range savedMessages {
		fullMessage, err := h.service.Message.GetFullMessageByMessageID(msg.MessageID)
		if err != nil {
			h.log.Error(err)
		}

		fullMessage.SavedID = msg.ID

		messages = append(messages, *fullMessage)
	}

	webUserID, webUserEmail := util.ProcessWebUserData(user)

	data := SavedPageData{
		DefaultPageData: PageData{
			Type:           "saved",
			Title:          "Saved user messages",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   webUserEmail,
			WebUserID:      webUserID,
		},
		Messages:       messages,
		MessagesLength: len(messages),
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) createSavedMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
	messageID, _ := strconv.Atoi(mux.Vars(r)["message_id"])

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	err = h.service.Saved.CreateSavedMessage(&model.Saved{WebUserID: userID, MessageID: messageID})
	if err != nil {
		h.log.Error(err)

		http.Redirect(w, r, "/home", http.StatusConflict)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%d", user.ID), http.StatusFound)
}

func (h *Handler) deleteSavedMessage(w http.ResponseWriter, r *http.Request) {
	savedID, _ := strconv.Atoi(mux.Vars(r)["saved_id"])

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	err = h.service.Saved.DeleteSavedMessage(savedID)
	if err != nil {
		h.log.Error(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%d", user.ID), http.StatusMovedPermanently)
}
