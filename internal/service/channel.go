package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type ChannelDBService struct {
	store *store.Store
}

func NewChannelDBService(store *store.Store) *ChannelDBService {
	return &ChannelDBService{
		store: store,
	}
}

func (s *ChannelDBService) CreateChannel(channel *model.DBChannel) error {
	err := s.store.Channel.CreateChannel(channel)
	if err != nil {
		return fmt.Errorf("[Channel] Service.CreateChannel error: %w", err)
	}

	return nil
}

func (s *ChannelDBService) GetChannels() ([]model.Channel, error) {
	channels, err := s.store.Channel.GetChannels()
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannels error: %w", err)
	}

	return channels, nil
}

func (s *ChannelDBService) GetChannelsByPage(page int) ([]model.Channel, error) {
	if page == 1 || page == 0 {
		page = 0
	} else if page == 2 {
		page = 10
	} else {
		page *= 10
		page -= 10
	}

	channels, err := s.store.Channel.GetChannelsByPage(page)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelsByPage error: %w", err)
	}

	return channels, nil
}

func (s *ChannelDBService) GetChannelByName(name string) (*model.Channel, error) {
	channel, err := s.store.Channel.GetChannelByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelByName error: %w", err)
	}

	return channel, nil
}

func (s *ChannelDBService) GetChannelStats(channelID int) (*model.Stat, error) {
	stat, err := s.store.Channel.GetChannelStats(channelID)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelStats error: %w", err)
	}

	return stat, nil
}
