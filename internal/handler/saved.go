package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/gorilla/mux"
)

type SavedPageData struct {
	Type           string
	Title          string
	UserEmail      string
	WebUserID      int
	Messages       []model.FullMessage
	MessagesLength int
	Channels       []model.Channel
	ChannelsLength int
}

func (h *Handler) savedPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["user_id"]
	userID, _ := strconv.Atoi(ID)

	data := SavedPageData{
		Title: "Saved user messages",
		Type:  "saved",
	}

	fullMessages := make([]model.FullMessage, 0)

	savedMessages, err := h.service.Saved.GetSavedMessages(userID)
	if err != nil {
		h.log.Error(err)
	}

	channels, err := h.service.Channel.GetChannels()
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

		fullMessages = append(fullMessages, *fullMessage)
	}

	data.Channels = util.ProcessChannels(channels)
	data.ChannelsLength = len(channels)
	data.Messages = fullMessages
	data.MessagesLength = len(fullMessages)
	data.WebUserID, data.UserEmail = util.ProcessWebUserData(user)

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) deleteSavedMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["saved_id"]

	savedID, _ := strconv.Atoi(ID)

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	err = h.service.Saved.DeleteSavedMessage(savedID)
	if err != nil {
		h.log.Error(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%d", user.ID), http.StatusFound)
}

func (h *Handler) createSavedMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uID := vars["user_id"]
	mID := vars["message_id"]

	userID, _ := strconv.Atoi(uID)
	messageID, _ := strconv.Atoi(mID)

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	err = h.service.Saved.CreateSavedMessage(&model.Saved{WebUserID: userID, MessageID: messageID})
	if err != nil {
		h.log.Error(err)

		http.Redirect(w, r, "/home", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/saved/%d", user.ID), http.StatusFound)
}
