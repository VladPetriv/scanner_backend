package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/convert"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type channelService struct {
	store   *store.Store
	logger  *logger.Logger
	message MessageService
}

func NewChannelService(store *store.Store, logger *logger.Logger, messageService MessageService) ChannelService {
	return &channelService{
		store:   store,
		logger:  logger,
		message: messageService,
	}
}

func (s channelService) CreateChannel(channel *model.DBChannel) error {
	logger := s.logger

	candidate, err := s.GetChannelByName(channel.Name)
	if err != nil {
		if !errors.Is(err, ErrChannelNotFound) {
			logger.Error().Err(err).Msg("get channel by name")
			return err
		}
	}
	if candidate != nil {
		logger.Info().Msg("channel exists")
		return ErrChannelExists
	}

	err = s.store.Channel.CreateChannel(channel)
	if err != nil {
		logger.Error().Err(err).Msg("create channel")
		return fmt.Errorf("create channel error: %w", err)
	}

	logger.Info().Msg("channel successfully created")
	return nil
}

func (s channelService) GetChannels() ([]model.Channel, error) {
	logger := s.logger

	channels, err := s.store.Channel.GetChannels()
	if err != nil {
		logger.Error().Err(err).Msg("get channels")
		return nil, fmt.Errorf("get channels error: %w", err)
	}

	if channels == nil {
		logger.Info().Msg("channels not found")
		return nil, ErrChannelsNotFound
	}

	logger.Info().Interface("channels", channels).Msg("successfully got channels")
	return channels, nil
}

func (s channelService) GetChannelsByPage(page int) ([]model.Channel, error) {
	logger := s.logger

	channels, err := s.store.Channel.GetChannelsByPage(convert.PageToOffset(page))
	if err != nil {
		logger.Error().Err(err).Msg("get channels by page")
		return nil, fmt.Errorf("get channels by page error: %w", err)
	}

	if channels == nil {
		logger.Info().Msg("channels by page not found")
		return nil, ErrChannelsNotFound
	}

	logger.Info().Interface("channels", channels).Msg("successfully got channels by page")
	return channels, nil
}

func (s channelService) GetChannelByName(name string) (*model.Channel, error) {
	logger := s.logger

	channel, err := s.store.Channel.GetChannelByName(name)
	if err != nil {
		logger.Error().Err(err).Msg("get channel by name")
		return nil, fmt.Errorf("get channel by name error: %w", err)
	}

	if channel == nil {
		logger.Info().Msg("channel not found")
		return nil, ErrChannelNotFound
	}

	logger.Info().Interface("channel", channel).Msg("successfully got channel by name")
	return channel, nil
}

func (s channelService) GetChannelStats(channelID int) (*model.Stat, error) {
	logger := s.logger

	stat, err := s.store.Channel.GetChannelStats(channelID)
	if err != nil {
		logger.Error().Err(err).Msg("get channels statistic")
		return nil, fmt.Errorf("get channels statistic error: %w", err)
	}

	logger.Info().Interface("statistic", stat).Msg("successfully got channel statistic")
	return stat, nil
}

func (s channelService) ProcessChannelPage(channelName string, page int) (*LoadChannelOutput, error) {
	logger := s.logger

	channel, err := s.GetChannelByName(channelName)
	if err != nil {
		if errors.Is(err, ErrChannelNotFound) {
			logger.Info().Msg("channel not found")
			return nil, err
		}

		logger.Error().Err(err).Msg("get channel by name")
		return nil, fmt.Errorf("get channel by name: %w", err)
	}

	messagesCount, err := s.message.GetMessagesCountByChannelID(channel.ID)
	if err != nil {
		if !errors.Is(err, ErrMessagesCountNotFound) {
			logger.Error().Err(err).Msg("get messages count by channel id")
			return nil, fmt.Errorf("get message count by channel id: %w", err)
		}
	}

	messages, err := s.message.GetFullMessagesByChannelIDAndPage(channel.ID, page)
	if err != nil {
		if !errors.Is(err, ErrMessagesNotFound) {
			logger.Error().Err(err).Msg("get messages by channel id and page")
			return nil, fmt.Errorf("get messages by channel id and page")
		}
	}

	return &LoadChannelOutput{
		Channel:       *channel,
		MessagesCount: messagesCount,
		Messages:      messages,
	}, nil
}

func (s channelService) ProcessChannelsPage(page int) (*LoadChannelsOutput, error) {
	logger := s.logger

	channels, err := s.GetChannelsByPage(page)
	if err != nil {
		if !errors.Is(err, ErrChannelsNotFound) {
			logger.Error().Err(err).Msg("get channels by page")
			return nil, fmt.Errorf("get channels by page: %w", err)
		}
	}

	for index, channel := range channels {
		stat, err := s.GetChannelStats(channel.ID)
		if err != nil {
			logger.Error().Err(err).Msg("get channel stats")

			continue
		}

		channels[index].Stats = *stat
	}

	logger.Info().Interface("channels", channels).Msg("successfully got channels with statistic")
	return &LoadChannelsOutput{
		Channels: channels,
	}, nil
}
