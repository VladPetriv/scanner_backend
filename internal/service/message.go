package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type MessageDBService struct {
	store *store.Store
}

func NewMessageDBService(store *store.Store) *MessageDBService {
	return &MessageDBService{
		store: store,
	}
}

func (s *MessageDBService) GetMessagesLength() (int, error) {
	length, err := s.store.Message.GetMessagesLength()
	if err != nil {
		return 0, fmt.Errorf("[Message] Service.GetMessagesLength error: %w", err)
	}

	if length == 0 {
		return 0, fmt.Errorf("messages not found")
	}

	return length, nil
}

func (s *MessageDBService) GetFullMessages(page int) ([]model.FullMessage, error) {
	if page == 1 || page == 0 {
		page = 0
	} else if page != 0 {
		page *= 10
		page -= 10
	}

	messages, err := s.store.Message.GetFullMessages(page)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessages error: %w", err)
	}

	if messages == nil {
		return nil, fmt.Errorf("messages not found")
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessagesByChannelID(ID, limit, page int) ([]model.FullMessage, error) {
	if page == 1 || page == 0 {
		page = 0
	} else if page != 0 {
		page *= 10
		page -= 10
	}

	messages, err := s.store.Message.GetFullMessagesByChannelID(ID, limit, page)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByChannelID error: %w", err)
	}

	if messages == nil {
		return nil, fmt.Errorf("messages not found")
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByUserID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByUserID error: %w", err)
	}

	if messages == nil {
		return nil, fmt.Errorf("messages not found")
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessageByMessageID(ID int) (*model.FullMessage, error) {
	message, err := s.store.Message.GetFullMessageByMessageID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessageByMessageID error: %w", err)
	}

	if message == nil {
		return nil, fmt.Errorf("messages not found")
	}

	return message, nil
}

func (s *MessageDBService) GetMessagesLengthByChannelID(ID int) (int, error) {
	length, err := s.store.Message.GetMessagesLengthByChannelID(ID)
	if err != nil {
		return 0, fmt.Errorf("[Message] Service.GetMessagesLengthByChannelID error: %w", err)
	}

	if length == 0 {
		return 0, fmt.Errorf("messages not found")
	}

	return length, nil
}
