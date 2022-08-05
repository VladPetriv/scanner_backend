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

func (s *MessageDBService) CreateMessage(message *model.DBMessage) (int, error) {
	id, err := s.store.Message.CreateMessage(message)
	if err != nil {
		return id, fmt.Errorf("[Message] Service.CreateMessage error: %w", err)
	}

	return id, nil
}

func (s *MessageDBService) GetMessagesCount() (int, error) {
	count, err := s.store.Message.GetMessagesCount()
	if err != nil {
		return 0, fmt.Errorf("[Message] Service.GetMessagesCount error: %w", err)
	}

	return count, nil
}

func (s *MessageDBService) GetMessagesCountByChannelID(ID int) (int, error) {
	count, err := s.store.Message.GetMessagesCountByChannelID(ID)
	if err != nil {
		return 0, fmt.Errorf("[Message] Service.GetMessagesCountByChannelID error: %w", err)
	}

	return count, nil
}

func (s *MessageDBService) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	if page == 1 || page == 0 {
		page = 0
	} else if page != 0 {
		page *= 10
		page -= 10
	}

	messages, err := s.store.Message.GetFullMessagesByPage(page)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByPage error: %w", err)
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessagesByChannelIDAndPage(ID, page int) ([]model.FullMessage, error) {
	if page == 1 || page == 0 {
		page = 0
	} else if page != 0 {
		page *= 10
		page -= 10
	}

	messages, err := s.store.Message.GetFullMessagesByChannelIDAndPage(ID, page)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByChannelIDAndPage error: %w", err)
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByUserID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByUserID error: %w", err)
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessageByMessageID(ID int) (*model.FullMessage, error) {
	message, err := s.store.Message.GetFullMessageByMessageID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessageByMessageID error: %w", err)
	}

	return message, nil
}
