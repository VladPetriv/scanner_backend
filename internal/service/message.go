package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/convert"
)

type messageService struct {
	store *store.Store
}

func NewMessageService(store *store.Store) MessageService {
	return &messageService{
		store: store,
	}
}

func (s messageService) CreateMessage(message *model.DBMessage) (int, error) {
	id, err := s.store.Message.CreateMessage(message)
	if err != nil {
		return id, fmt.Errorf("[Message] Service.CreateMessage error: %w", err)
	}

	return id, nil
}

func (s messageService) GetMessagesCount() (int, error) {
	count, err := s.store.Message.GetMessagesCount()
	if err != nil {
		return 0, fmt.Errorf("[Message] Service.GetMessagesCount error: %w", err)
	}

	if count == 0 {
		return 0, ErrMessagesCountNotFound
	}

	return count, nil
}

func (s messageService) GetMessagesCountByChannelID(id int) (int, error) {
	count, err := s.store.Message.GetMessagesCountByChannelID(id)
	if err != nil {
		return 0, fmt.Errorf("[Message] Service.GetMessagesCountByChannelID error: %w", err)
	}

	if count == 0 {
		return 0, ErrMessagesCountNotFound
	}

	return count, nil
}

func (s messageService) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByPage(convert.PageToOffset(page))
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByPage error: %w", err)
	}

	if messages == nil {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (s messageService) GetFullMessagesByChannelIDAndPage(id, page int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByChannelIDAndPage(id, convert.PageToOffset(page))
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByChannelIDAndPage error: %w", err)
	}

	if messages == nil {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (s messageService) GetFullMessagesByUserID(id int) ([]model.FullMessage, error) {
	messages, err := s.store.Message.GetFullMessagesByUserID(id)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessagesByUserID error: %w", err)
	}

	if messages == nil {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (s messageService) GetFullMessageByMessageID(id int) (*model.FullMessage, error) {
	message, err := s.store.Message.GetFullMessageByID(id)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetFullMessageByMessageID error: %w", err)
	}

	if message == nil {
		return nil, ErrMessageNotFound
	}

	return message, nil
}
