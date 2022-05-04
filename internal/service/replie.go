package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type ReplieDBService struct {
	store *store.Store
}

func NewReplieDBService(store *store.Store) *ReplieDBService {
	return &ReplieDBService{store: store}
}

func (s *ReplieDBService) GetReplies() ([]model.Replie, error) {
	replies, err := s.store.Replie.GetReplies()
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplies error: %w", err)
	}

	if replies == nil {
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
}

func (s *ReplieDBService) GetReplie(replieID int) (*model.Replie, error) {
	replie, err := s.store.Replie.GetReplie(replieID)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplie error: %w", err)
	}

	if replie == nil {
		return nil, fmt.Errorf("replie not found")
	}

	return replie, nil
}

func (s *ReplieDBService) GetReplieByName(name string) (*model.Replie, error) {
	replie, err := s.store.Replie.GetReplieByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplieByName error: %w", err)
	}

	if replie == nil {
		return nil, nil
	}

	return replie, nil
}

func (s *ReplieDBService) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies, err := s.store.Replie.GetFullRepliesByMessageID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetFullRepliesByMessageID error: %w", err)
	}

	if len(replies) == 0 {
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
}
