package handler

import (
	"net/http"
	"strconv"

	"github.com/VladPetriv/go-pagination-bootstrap"
	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

const channelsPerPage = 10

type channelPageData struct {
	DefaultPageData PageData
	Channel         model.Channel
	Messages        []model.FullMessage
	MessagesLength  int
	Pager           *pagination.Pagination
}

type channelsPageData struct {
	DefaultPageData PageData
	Channels        []model.Channel
	Pager           *pagination.Pagination
}

func (h Handler) loadChannelsPage(w http.ResponseWriter, r *http.Request) {
	data := channelsPageData{
		DefaultPageData: PageData{
			Title:        "Telegram channels",
			Type:         "channels",
			WebUserEmail: "",
			WebUserID:    0,
		},
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error().Err(err).Msg("convert page value for channels to int")
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

	pageData, err := h.service.Channel.ProcessChannelsPage(page)
	if err != nil {
		h.log.Error().Err(err).Msg("get data for channels page")
	}
	if pageData != nil {
		data.Channels = pageData.Channels
		data.Pager = pagination.New(len(navBarChannels), channelsPerPage, page, "/channel/?page=2")
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load channels page")
	}
}

func (h Handler) loadChannelPage(w http.ResponseWriter, r *http.Request) {
	data := channelPageData{
		DefaultPageData: PageData{
			Type:         "channel",
			Title:        "Telegram channel",
			WebUserEmail: "",
			WebUserID:    0,
		},
	}

	channelName := mux.Vars(r)["channel_name"]

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error().Err(err).Msg("convert page value for channel to int")
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

	pageData, err := h.service.Channel.ProcessChannelPage(channelName, page)
	if err != nil {
		h.log.Error().Err(err).Msg("get data for channel page")
	}
	if pageData != nil {
		data.Channel = pageData.Channel
		data.Messages = pageData.Messages
		data.MessagesLength = pageData.MessagesCount
		data.Pager = pagination.New(len(pageData.Messages), messagesPerPage, page, "/channel/ru_python/?page=0")
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("load channel page")
	}
}
