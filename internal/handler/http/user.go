package handler

import (
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

func (h Handler) loadUserPage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("covert user id to int")

		http.Redirect(w, r, "/home", http.StatusNotFound)
		return
	}

	tgUser, err := h.service.User.GetUserByID(userID)
	if err != nil {
		h.log.Error().Err(err).Msg("get user by id")
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for navbar")
	}

	messages, err := h.service.Message.GetFullMessagesByUserID(tgUser.ID)
	if err != nil {
		h.log.Error().Err(err).Msg("get full messages by user id")
	}

	user, err := h.service.WebUser.GetWebUserByEmail(h.getUserFromSession(r))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	data := UserPageData{
		DefaultPageData: PageData{
			Type:           "user",
			Title:          "Telegram User",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   "",
			WebUserID:      0,
		},
		User:           *tgUser,
		Messages:       messages,
		MessagesLength: len(messages),
	}
	if user != nil {
		data.DefaultPageData.WebUserEmail = user.Email
		data.DefaultPageData.WebUserID = user.ID
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load user page")
	}
}
