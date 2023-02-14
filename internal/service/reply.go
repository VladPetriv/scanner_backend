package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type replyService struct {
	store  *store.Store
	logger *logger.Logger
}

var _ ReplyService = (*replyService)(nil)

func NewReplyService(store *store.Store, logger *logger.Logger) *replyService {
	return &replyService{
		store:  store,
		logger: logger,
	}
}

func (s replyService) CreateReply(reply *model.DBReply) error {
	logger := s.logger

	err := s.store.Reply.CreateReply(reply)
	if err != nil {
		logger.Error().Err(err).Msg("create reply")
		return fmt.Errorf("create reply error: %w", err)
	}

	logger.Info().Msg("reply successfully created")
	return nil
}

func (s replyService) GetFullRepliesByMessageID(messageID int) ([]model.FullReply, error) {
	logger := s.logger

	replies, err := s.store.Reply.GetFullRepliesByMessageID(messageID)
	if err != nil {
		logger.Error().Err(err).Msg("get replies by message id")
		return nil, fmt.Errorf("get replies by message id error: %w", err)
	}

	if replies == nil {
		logger.Info().Msg("replies not found")
		return nil, ErrRepliesNotFound
	}

	logger.Info().Interface("replies", replies).Msg("successfully got replies by message id")
	return replies, nil
}
