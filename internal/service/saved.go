package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type savedService struct {
	store   *store.Store
	logger  *logger.Logger
	message MessageService
}

var _ SavedService = (*savedService)(nil)

func NewSavedService(store *store.Store, logger *logger.Logger, messageService MessageService) *savedService {
	return &savedService{
		store:   store,
		logger:  logger,
		message: messageService,
	}
}

func (s savedService) GetSavedMessageByMessageID(id int) (*model.Saved, error) {
	logger := s.logger

	message, err := s.store.Saved.GetSavedMessageByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("get saved message by message id")
		return nil, fmt.Errorf("get saved message by message id from db: %w", err)
	}
	if message == nil {
		logger.Info().Msg("saved message not found")
		return nil, ErrSavedMessageNotFound
	}

	logger.Info().Interface("message", message).Msg("successfully got saved message by message id")
	return message, nil
}

func (s savedService) CreateSavedMessage(savedMessage *model.Saved) error {
	logger := s.logger

	err := s.store.Saved.CreateSavedMessage(savedMessage)
	if err != nil {
		logger.Error().Err(err).Msg("create saved message")
		return fmt.Errorf("create saved message in db: %w", err)
	}

	logger.Info().Msg("saved message successfully created")
	return nil
}

func (s savedService) DeleteSavedMessage(id int) error {
	logger := s.logger

	err := s.store.Saved.DeleteSavedMessage(id)
	if err != nil {
		logger.Error().Err(err).Msg("delete saved message")
		return fmt.Errorf("delete saved message from db: %w", err)
	}

	logger.Info().Msg("successfully delete saved message")
	return nil
}

func (s savedService) ProcessSavedMessages(userID int) (*LoadSavedMessagesOutput, error) {
	logger := s.logger

	var savedMessages []model.FullMessage

	messages, err := s.store.Saved.GetSavedMessages(userID)
	if err != nil {
		if errors.Is(err, ErrSavedMessagesNotFound) {
			logger.Info().Msg("messages not found")
			return &LoadSavedMessagesOutput{}, nil
		}

		logger.Error().Err(err).Msg("get saved messages")
		return nil, fmt.Errorf("get saved messages from db: %w", err)
	}

	for _, msg := range messages {
		fullMessage, err := s.message.GetFullMessageByMessageID(msg.MessageID)
		if err != nil {
			if errors.Is(err, ErrMessageNotFound) {
				logger.Info().Msg("message not found")
			} else {
				logger.Error().Err(err).Msg("get full messages by message id")
			}

			continue
		}

		fullMessage.SavedID = msg.ID

		savedMessages = append(savedMessages, *fullMessage)
	}

	return &LoadSavedMessagesOutput{
		SavedMessages:      savedMessages,
		SavedMessagesCount: len(savedMessages),
	}, nil
}
