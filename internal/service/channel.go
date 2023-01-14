package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/convert"
)

type channelService struct {
	store *store.Store
}

func NewChannelService(store *store.Store) ChannelService {
	return &channelService{
		store: store,
	}
}

func (s channelService) CreateChannel(channel *model.DBChannel) error {
	candidate, err := s.GetChannelByName(channel.Name)
	if err != nil {
		if !errors.Is(err, ErrChannelNotFound) {
			return err
		}
	}
	if candidate != nil {
		return ErrChannelExists
	}

	err = s.store.Channel.CreateChannel(channel)
	if err != nil {
		return fmt.Errorf("[Channel] Service.CreateChannel error: %w", err)
	}

	return nil
}

func (s channelService) GetChannels() ([]model.Channel, error) {
	channels, err := s.store.Channel.GetChannels()
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannels error: %w", err)
	}

	if channels == nil {
		return nil, ErrChannelsNotFound
	}

	return channels, nil
}

func (s channelService) GetChannelsByPage(page int) ([]model.Channel, error) {
	channels, err := s.store.Channel.GetChannelsByPage(convert.PageToOffset(page))
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelsByPage error: %w", err)
	}

	if channels == nil {
		return nil, ErrChannelsNotFound
	}

	return channels, nil
}

func (s channelService) GetChannelByName(name string) (*model.Channel, error) {
	channel, err := s.store.Channel.GetChannelByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelByName error: %w", err)
	}

	if channel == nil {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}

func (s channelService) GetChannelStats(channelID int) (*model.Stat, error) {
	stat, err := s.store.Channel.GetChannelStats(channelID)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelStats error: %w", err)
	}

	return stat, nil
}
