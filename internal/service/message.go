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

var _ MessageService = (*messageService)(nil)

func NewMessageService(store *store.Store, logger *logger.Logger, replyService ReplyService) *messageService {
	return &messageService{
		store:  store,
		logger: logger,
		reply:  replyService,
	}
}

func (s messageService) CreateMessage(message *model.DBMessage) (int, error) {
	logger := s.logger

	candidate, err := s.store.Message.GetMessageByTitle(message.Title)
	if err != nil {
		logger.Error().Err(err).Msg("get message by title")
		return 0, fmt.Errorf("get message by title from db: %w", err)
	}
	if candidate == nil {
		logger.Info().Str("message title", message.Title).Msg("message not found")
	}

	if candidate != nil && candidate.ChannelID == message.ChannelID {
		logger.Info().Msg("message is exist")
		return 0, ErrMessageExists
	}

	id, err := s.store.Message.CreateMessage(message)
	if err != nil {
		logger.Error().Err(err).Msg("create message")
		return id, fmt.Errorf("create message in db: %w", err)
	}

	logger.Info().Int("message id", id).Msg("message successfully created")
	return id, nil
}

func (s messageService) GetMessagesCount() (int, error) {
	logger := s.logger

	count, err := s.store.Message.GetMessagesCount()
	if err != nil {
		logger.Error().Err(err).Msg("get messages count")
		return 0, fmt.Errorf("get messages count from db: %w", err)
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
		return 0, fmt.Errorf("get messages count by channel id from db: %w", err)
	}

	if count == 0 {
		logger.Info().Msg("messages count by channel id not found")
		return 0, ErrMessagesCountNotFound
	}

	logger.Info().Int("messages count", count).Msg("successfully got messages count by channel id")
	return count, nil
}

func (s messageService) GetFullMessagesByChannelIDAndPage(id, page int) ([]model.FullMessage, error) {
	logger := s.logger

	messages, err := s.store.Message.GetFullMessagesByChannelIDAndPage(id, convert.PageToOffset(page))
	if err != nil {
		logger.Error().Err(err).Msg("get full messages by channel id and page")
		return nil, fmt.Errorf("get full messages by channel id and page from db: %w", err)
	}

	if messages == nil {
		logger.Info().Msg("full messages by channel id and page not found")
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
		return nil, fmt.Errorf("get full messages by user id from db: %w", err)
	}

	if messages == nil {
		logger.Info().Msg("full messages by user id not found")
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
		return nil, fmt.Errorf("get full message by message id from db: %w", err)
	}

	if message == nil {
		logger.Info().Int("message id", id).Msg("message by id not found")
		return nil, ErrMessageNotFound
	}

	logger.Info().Interface("message", message).Msg("successfully got message by id")
	return message, nil
}

func (s messageService) ProcessMessagePage(messageID int) (*LoadMessageOutput, error) {
	logger := s.logger

	message, err := s.store.Message.GetFullMessageByID(messageID)
	if err != nil {
		logger.Error().Err(err).Msg("get full message by id")
		return nil, fmt.Errorf("get full message by id from db: %w", err)
	}
	if message == nil {
		logger.Info().Int("message id", messageID).Msg("message by id not found")
		return &LoadMessageOutput{}, nil
	}

	replies, err := s.reply.GetFullRepliesByMessageID(message.ID)
	if err != nil {
		if errors.Is(err, ErrRepliesNotFound) {
			logger.Info().Int("message id", message.ID).Msg("full replies by message id not found")
			return &LoadMessageOutput{
				Message: message,
			}, nil
		}

		logger.Error().Err(err).Msg("get replies by message id")
		return nil, fmt.Errorf("[ProcessMessagePage]: %w", err)
	}

	message.Replies = replies
	message.RepliesCount = len(replies)

	return &LoadMessageOutput{Message: message}, nil
}

func (s messageService) ProcessHomePage(page int) (*LoadHomeOutput, error) {
	logger := s.logger

	messagesCount, err := s.store.Message.GetMessagesCount()
	if err != nil {
		logger.Error().Err(err).Msg("get messages count")
		return nil, fmt.Errorf("get messages count from db: %w", err)
	}
	if messagesCount == 0 {
		logger.Info().Msg("message count not found")
		return &LoadHomeOutput{}, nil
	}

	messages, err := s.store.Message.GetFullMessagesByPage(convert.PageToOffset(page))
	if err != nil {
		logger.Error().Err(err).Msg("get messages by page")
		return nil, fmt.Errorf("get full messages by page from db: %w", err)
	}

	return &LoadHomeOutput{
		Messages:      messages,
		MessagesCount: messagesCount,
	}, nil
}
