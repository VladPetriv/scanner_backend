package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/gorilla/mux"
)

type UserPageData struct {
	Type           string
	Title          string
	User           model.User
	Channels       []model.Channel
	ChannelsLength int
	Messages       []model.FullMessage
	MessagesLength int
	WebUserID      int
	UserEmail      string
}

func (h *Handler) userPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	data := UserPageData{
		Type:  "user",
		Title: "Telegram User",
	}

	ID, _ := strconv.Atoi(userID)

	user, err := h.service.User.GetUserByID(ID)
	if err != nil {
		h.log.Error(err)
	}

	channels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error(err)
	}

	messages, err := h.service.Message.GetFullMessagesByUserID(user.ID)
	if err != nil {
		h.log.Error(err)
	}

	webUser, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error(err)
	}

	data.User = *user
	data.Channels = util.ProcessChannels(channels)
	data.ChannelsLength = len(channels)
	data.Messages = messages
	data.MessagesLength = len(messages)
	data.WebUserID, data.UserEmail = util.ProcessWebUserData(webUser)

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error(err)
	}
}
