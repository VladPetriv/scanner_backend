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

func (s *ReplieDBService) CreateReplie(replie *model.DBReplie) error {
	err := s.store.Replie.CreateReplie(replie)
	if err != nil {
		return fmt.Errorf("[Replie] Service.CreateReplie error: %w", err)
	}

	return nil
}

func (s *ReplieDBService) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies, err := s.store.Replie.GetFullRepliesByMessageID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetFullRepliesByMessageID error: %w", err)
	}

	return replies, nil
}
