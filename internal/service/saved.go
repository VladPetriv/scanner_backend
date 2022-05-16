package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type SavedDbService struct {
	store *store.Store
}

func NewSavedDbService(store *store.Store) *SavedDbService {
	return &SavedDbService{store: store}
}

func (s *SavedDbService) GetSavedMessages(UserID int) ([]model.Saved, error) {
	savedMessages, err := s.store.Saved.GetSavedMessages(UserID)
	if err != nil {
		return nil, fmt.Errorf("[Saved] Service.GetSavedMessages error: %w", err)
	}

	if savedMessages == nil {
		return nil, fmt.Errorf("saved messages not found")
	}

	return savedMessages, nil
}

func (s *SavedDbService) GetSavedMessageByMessageID(ID int) (*model.Saved, error) {
	savedMessage, err := s.store.Saved.GetSavedMessageByMessageID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Saved] Service.GetSavedMessageByMessageID error: %w", err)
	}

	if savedMessage == nil {
		return nil, fmt.Errorf("saved message not found")
	}

	return savedMessage, nil
}

func (s *SavedDbService) CreateSavedMessage(savedMessage *model.Saved) error {
	_, err := s.store.Saved.CreateSavedMessage(savedMessage)
	if err != nil {
		return fmt.Errorf("[Saved] Service.CreateSaveMessage error: %w", err)
	}

	return nil
}

func (s *SavedDbService) DeleteSavedMessage(ID int) error {
	_, err := s.store.Saved.DeleteSavedMessage(ID)
	if err != nil {
		return fmt.Errorf("[Saved] Service.DeleteSavedMessage error: %w", err)
	}

	return nil
}
