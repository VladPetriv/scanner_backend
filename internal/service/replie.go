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
