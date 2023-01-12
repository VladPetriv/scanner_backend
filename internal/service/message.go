package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/convert"
)

var (
	ErrMessagesCountNotFound = errors.New("message count not found")
	ErrMessagesNotFound      = errors.New("messages not found")
	ErrMessageNotFound       = errors.New("messages not found")
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

	if count == 0 {
		return 0, ErrMessagesCountNotFound
	}

	return count, nil
}

func (s *MessageDBService) GetMessagesCountByChannelID(id int) (int, error) {
	count, err := s.store.Message.GetMessagesCountByChannelID(id)
	if err != nil {
		return 0, fmt.Errorf("[Message] Service.GetMessagesCountByChannelID error: %w", err)
	}

	if count == 0 {
		return 0, ErrMessagesCountNotFound
	}

	return count, nil
}

func (s *MessageDBService) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByPage(convert.PageToOffset(page))
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByPage error: %w", err)
	}

	if messages == nil {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessagesByChannelIDAndPage(id, page int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByChannelIDAndPage(id, convert.PageToOffset(page))
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByChannelIDAndPage error: %w", err)
	}

	if messages == nil {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessagesByUserID(id int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByUserID(id)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByUserID error: %w", err)
	}

	if messages == nil {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (s *MessageDBService) GetFullMessageByMessageID(id int) (*model.FullMessage, error) {
	message, err := s.store.Message.GetFullMessageByID(id)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessageByMessageID error: %w", err)
	}

	if message == nil {
		return nil, ErrMessageNotFound
	}

	return message, nil
}
