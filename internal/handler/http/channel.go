package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladPetriv/go-pagination-bootstrap"
	"github.com/gorilla/mux"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

type ChannelPageData struct {
	DefaultPageData PageData
	Channel         model.Channel
	Channels        []model.Channel
	Messages        []model.FullMessage
	MessagesLength  int
	Pager           *pagination.Pagination
}

func (h Handler) channelsPage(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error().Err(err).Msg("convert page value to int")
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for navbar")
	}

	channels, err := h.service.Channel.GetChannelsByPage(page)
	if err != nil {
		h.log.Error().Err(err).Msg("get channels by page")
	}

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	for index, channel := range channels {
		stat, err := h.service.Channel.GetChannelStats(channel.ID)
		if err != nil {
			h.log.Error().Err(err).Msg("get channel stats")
		}

		channels[index].Stats = *stat
	}

	webUserID, webUserEmail := util.ProcessWebUserData(user)

	data := ChannelPageData{
		DefaultPageData: PageData{
			Title:          "Telegram channels",
			Type:           "channels",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   webUserEmail,
			WebUserID:      webUserID,
		},
		Channels: channels,
		Pager:    pagination.New(len(navBarChannels), 10, page, "/channel/?page=2"),
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("execute base template")
	}
}

func (h Handler) channelPage(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["channel_name"]
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error().Err(err).Msg("convert page to int")
	}

	navBarChannels, err := h.service.Channel.GetChannels()
	if err != nil {
		h.log.Error().Err(err).Msg("get channels for navbar")
	}

	channel, err := h.service.Channel.GetChannelByName(name)
	if err != nil {
		h.log.Error().Err(err).Msg("get channel by name")
	}

	count, err := h.service.Message.GetMessagesCountByChannelID(channel.ID)
	if err != nil {
		h.log.Error().Err(err).Msg("get messages count by id")
	}

	messages, err := h.service.Message.GetFullMessagesByChannelIDAndPage(channel.ID, page)
	if err != nil {
		h.log.Error().Err(err).Msg("get full messages by channel id and page")
	}

	user, err := h.service.WebUser.GetWebUserByEmail(fmt.Sprint(h.checkUserStatus(r)))
	if err != nil {
		h.log.Error().Err(err).Msg("get web user by email")
	}

	webUserID, webUserEmail := util.ProcessWebUserData(user)

	data := ChannelPageData{
		DefaultPageData: PageData{
			Type:           "channel",
			Title:          "Telegram channel",
			Channels:       util.ProcessChannels(navBarChannels),
			ChannelsLength: len(navBarChannels),
			WebUserEmail:   webUserEmail,
			WebUserID:      webUserID,
		},
		Channel:        *channel,
		Messages:       messages,
		MessagesLength: count,
		Pager:          pagination.New(count, 10, page, "/channel/ru_python/?page=0"),
	}

	err = h.templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.log.Error().Err(err).Msg("execute base template")
	}
}
