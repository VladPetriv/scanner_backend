package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type UserPageData struct {
	DefaultPageData PageData
	User            model.User
	Messages        []model.FullMessage
	MessagesLength  int
}

func (h *Handler) userPage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("covert user id to int")
	}

	user, err := h.service.User.GetUserByID(userID)
	if err != nil {
		h.log.Error().Err(err).Msg("get user by id")
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for navbar")
	}

	messages, err := h.service.Message.GetFullMessagesByUserID(user.ID)
	if err != nil {
		h.log.Error().Err(err).Msg("get full messages by user id")
	}

	webUser, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	webUserID, webUserEmail := util.ProcessWebUserData(webUser)

	data := UserPageData{
		DefaultPageData: PageData{
			Type:           "user",
			Title:          "Telegram User",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   webUserEmail,
			WebUserID:      webUserID,
		},
		User:           *user,
		Messages:       messages,
		MessagesLength: len(messages),
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("execute user page")
	}
}
