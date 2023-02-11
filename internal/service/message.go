package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/convert"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type messageService struct {
	store  *store.Store
	logger *logger.Logger
	reply  ReplyService
}

func NewMessageService(store *store.Store, logger *logger.Logger, replyService ReplyService) MessageService {
	return &messageService{
		store:  store,
		logger: logger,
		reply:  replyService,
	}
}

func (s messageService) CreateMessage(message *model.DBMessage) (int, error) {
	logger := s.logger

	candidate, err := s.GetMessageByTitle(message.Title)
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			logger.Info().Msg("message not found")
			return 0, err
		}

		logger.Error().Err(err).Msg("get message by title")
		return 0, err
	}

	if candidate != nil && candidate.ChannelID == message.ChannelID {
		logger.Info().Msg("message is exist")
		return 0, ErrMessageExists
	}

	id, err := s.store.Message.CreateMessage(message)
	if err != nil {
		logger.Error().Err(err).Msg("create message")
		return id, fmt.Errorf("create message error: %w", err)
	}

	logger.Info().Int("message id", id).Msg("message successfully created")
	return id, nil
}

func (s messageService) GetMessagesCount() (int, error) {
	logger := s.logger

	count, err := s.store.Message.GetMessagesCount()
	if err != nil {
		logger.Error().Err(err).Msg("get messages count")
		return 0, fmt.Errorf("get messages count error: %w", err)
	}

	if count == 0 {
		logger.Info().Msg("message count not found")
		return 0, ErrMessagesCountNotFound
	}

	logger.Info().Int("messages count", count).Msg("successfully got messages count")
	return count, nil
}

func (s messageService) GetMessagesCountByChannelID(id int) (int, error) {
	logger := s.logger

	count, err := s.store.Message.GetMessagesCountByChannelID(id)
	if err != nil {
		logger.Error().Err(err).Msg("get messages count by channel id")
		return 0, fmt.Errorf("get messages count by channel id error: %w", err)
	}

	if count == 0 {
		logger.Info().Msg("messages count by channel id not found")
		return 0, ErrMessagesCountNotFound
	}

	logger.Info().Int("messages count", count).Msg("successfully got messages count by channel id")
	return count, nil
}

func (s messageService) GetMessageByTitle(title string) (*model.DBMessage, error) {
	logger := s.logger

	message, err := s.store.Message.GetMessageByTitle(title)
	if err != nil {
		logger.Error().Err(err).Msg("get message by title")
		return nil, fmt.Errorf("get message by title error: %w", err)
	}

	if message == nil {
		logger.Info().Msg("message not found")
		return nil, ErrMessageNotFound
	}

	logger.Info().Interface("message", message).Msg("successfully got message by title")
	return message, nil
}

func (s messageService) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	logger := s.logger

	messages, err := s.store.Message.GetFullMessagesByPage(convert.PageToOffset(page))
	if err != nil {
		logger.Error().Err(err).Msg("get full messages by page")
		return nil, fmt.Errorf("get full messages by page error: %w", err)
	}

	if messages == nil {
		logger.Info().Msg("messages not found")
		return nil, ErrMessagesNotFound
	}

	logger.Info().Interface("messages", messages).Msg("successfully got message by page")
	return messages, nil
}

func (s messageService) GetFullMessagesByChannelIDAndPage(id, page int) ([]model.FullMessage, error) {
	logger := s.logger

	messages, err := s.store.Message.GetFullMessagesByChannelIDAndPage(id, convert.PageToOffset(page))
	if err != nil {
		logger.Error().Err(err).Msg("get full messages by page and channel id")
		return nil, fmt.Errorf("get full messages by page and channel id error: %w", err)
	}

	if messages == nil {
		logger.Info().Msg("messages not found")
		return nil, ErrMessagesNotFound
	}

	logger.Info().Interface("messages", messages).Msg("successfully got message by channel id and page")
	return messages, nil
}

func (s messageService) GetFullMessagesByUserID(id int) ([]model.FullMessage, error) {
	logger := s.logger

	messages, err := s.store.Message.GetFullMessagesByUserID(id)
	if err != nil {
		logger.Error().Err(err).Msg("get full messages by user id")
		return nil, fmt.Errorf("get full messages by user id error: %w", err)
	}

	if messages == nil {
		logger.Info().Msg("messages not found")
		return nil, ErrMessagesNotFound
	}

	logger.Info().Interface("messages", messages).Msg("successfully got full message by user id")
	return messages, nil
}

func (s messageService) GetFullMessageByMessageID(id int) (*model.FullMessage, error) {
	logger := s.logger

	message, err := s.store.Message.GetFullMessageByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("get full message by message id")
		return nil, fmt.Errorf("get full message by message id error: %w", err)
	}

	if message == nil {
		logger.Info().Msg("message not found")
		return nil, ErrMessageNotFound
	}

	logger.Info().Interface("message", message).Msg("successfully got message by id")
	return message, nil
}

func (s messageService) ProcessMessagePage(messageID int) (*LoadMessageOutput, error) {
	logger := s.logger

	message, err := s.GetFullMessageByMessageID(messageID)
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			logger.Info().Msg("message not found")
			return nil, err
		}

		logger.Error().Err(err).Msg("get message by id")
		return nil, fmt.Errorf("get message by id: %w", err)
	}

	replies, err := s.reply.GetFullRepliesByMessageID(message.ID)
	if err != nil {
		if errors.Is(err, ErrRepliesNotFound) {
			logger.Info().Msg("replies not found")
		}

		logger.Error().Err(err).Msg("get replies by message id")
		return nil, fmt.Errorf("get replies by message id: %w", err)
	}

	message.Replies = replies
	message.RepliesCount = len(replies)

	return &LoadMessageOutput{Message: message}, nil
}
