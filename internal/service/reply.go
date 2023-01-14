package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type ReplyDBService struct {
	store *store.Store
}

func NewReplyDBService(store *store.Store) *ReplyDBService {
	return &ReplyDBService{store: store}
}

func (s ReplyDBService) CreateReply(reply *model.DBReply) error {
	err := s.store.Reply.CreateReply(reply)
	if err != nil {
		return fmt.Errorf("[Reply] Service.CreateReply error: %w", err)
	}

	return nil
}

func (s ReplyDBService) GetFullRepliesByMessageID(messageID int) ([]model.FullReply, error) {
	replies, err := s.store.Reply.GetFullRepliesByMessageID(messageID)
	if err != nil {
		return nil, fmt.Errorf("[Reply] Service.GetFullRepliesByMessageID error: %w", err)
	}

	if replies == nil {
		//TODO: return custom error
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
}
