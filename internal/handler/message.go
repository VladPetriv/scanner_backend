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
	DefaultPageData PageData
	Message         model.FullMessage
}

func (h *Handler) messagePage(w http.ResponseWriter, r *http.Request) {
	messageID, _ := strconv.Atoi(mux.Vars(r)["message_id"])

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	message, err := h.service.Message.GetFullMessageByMessageID(messageID)
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

	webUserID, webUserEmail := util.ProcessWebUserData(user)

	data := MessagePageData{
		DefaultPageData: PageData{
			Type:           "message",
			Title:          "Telegram message",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   webUserEmail,
			WebUserID:      webUserID,
		},
		Message: *message,
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error(err)
	}
}
