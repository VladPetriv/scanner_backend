package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type replyService struct {
	store *store.Store
}

func NewReplyService(store *store.Store) ReplyService {
	return &replyService{store: store}
}

func (s replyService) CreateReply(reply *model.DBReply) error {
	err := s.store.Reply.CreateReply(reply)
	if err != nil {
		return fmt.Errorf("[Reply] Service.CreateReply error: %w", err)
	}

	return nil
}

func (s replyService) GetFullRepliesByMessageID(messageID int) ([]model.FullReply, error) {
	replies, err := s.store.Reply.GetFullRepliesByMessageID(messageID)
	if err != nil {
		return nil, fmt.Errorf("[Reply] Service.GetFullRepliesByMessageID error: %w", err)
	}

	if replies == nil {
		return nil, ErrRepliesNotFound
	}

	return replies, nil
}
