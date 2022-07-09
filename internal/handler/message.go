package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type MessagePageData struct {
	Type           string
	Title          string
	Channels       []model.Channel
	ChannelsLength int
	Message        model.FullMessage
	UserEmail      interface{}
	WebUserID      int
}

func (h *Handler) messagePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["message_id"]

	data := MessagePageData{
		Type:  "message",
		Title: "Telegram message",
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	ID, _ := strconv.Atoi(messageID)

	message, err := h.service.Message.GetFullMessageByMessageID(ID)
	if err != nil {
		h.log.Error(err)
	}

	replies, err := h.service.Replie.GetFullRepliesByMessageID(message.ID)
	if err != nil {
		h.log.Error(err)
	}

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	message.Replies = replies

	data.Channels = util.ProcessChannels(channels)
	data.ChannelsLength = len(channels)
	data.Message = *message
	data.WebUserID, data.UserEmail = util.ProcessWebUserData(user)

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error(err)
	}
}
