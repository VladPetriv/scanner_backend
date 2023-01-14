package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

var (
	ErrSavedMessagesNotFound = errors.New("saved messages not found")
	ErrSavedMessageNotFound  = errors.New("saved message not found")
)

type SavedDbService struct {
	store *store.Store
}

func NewSavedDbService(store *store.Store) *SavedDbService {
	return &SavedDbService{store: store}
}

func (s SavedDbService) GetSavedMessages(userID int) ([]model.Saved, error) {
	messages, err := s.store.Saved.GetSavedMessages(userID)
	if err != nil {
		return nil, fmt.Errorf("[Saved] Service.GetSavedMessages error: %w", err)
	}

	if messages == nil {
		return nil, ErrSavedMessagesNotFound
	}

	return messages, nil
}

func (s SavedDbService) GetSavedMessageByMessageID(id int) (*model.Saved, error) {
	message, err := s.store.Saved.GetSavedMessageByID(id)
	if err != nil {
		return nil, fmt.Errorf("[Saved] Service.GetSavedMessageByMessageID error: %w", err)
	}

	if message == nil {
		return nil, ErrSavedMessageNotFound
	}

	return message, nil
}

func (s SavedDbService) CreateSavedMessage(savedMessage *model.Saved) error {
	err := s.store.Saved.CreateSavedMessage(savedMessage)
	if err != nil {
		return fmt.Errorf("[Saved] Service.CreateSavedMessage error: %w", err)
	}

	return nil
}

func (s SavedDbService) DeleteSavedMessage(id int) error {
	err := s.store.Saved.DeleteSavedMessage(id)
	if err != nil {
		return fmt.Errorf("[Saved] Service.DeleteSavedMessage error: %w", err)
	}

	return nil
}
