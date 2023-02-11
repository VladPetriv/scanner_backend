package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type UserPageData struct {
	DefaultPageData PageData
	User            model.User
	Messages        []model.FullMessage
	MessagesLength  int
}

func (h Handler) loadUserPage(w http.ResponseWriter, r *http.Request) {
	data := UserPageData{
		DefaultPageData: PageData{
			Type:         "user",
			Title:        "Telegram User",
			WebUserEmail: "",
			WebUserID:    0,
		},
	}

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error().Err(err).Msg("covert user id to int")

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

	pageData, err := h.service.User.ProcessUserPage(userID)
	if err != nil {
		h.log.Error().Err(err).Msg("proccess user page")
	}
	if pageData != nil {
		data.Messages = pageData.Messages
		data.MessagesLength = len(pageData.Messages)
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load user page")
	}
}
