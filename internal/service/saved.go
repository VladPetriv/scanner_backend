package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type savedService struct {
	logger  *logger.Logger
	store   *store.Store
	message MessageService
}

func NewSavedService(store *store.Store, logger *logger.Logger, messageService MessageService) SavedService {
	return &savedService{
		logger:  logger,
		store:   store,
		message: messageService,
	}
}

func (s savedService) GetSavedMessages(userID int) ([]model.Saved, error) {
	logger := s.logger

	messages, err := s.store.Saved.GetSavedMessages(userID)
	if err != nil {
		logger.Error().Err(err).Msg("get saved message by user id")
		return nil, fmt.Errorf("get saved messages by user id: %w", err)
	}

	if messages == nil {
		logger.Info().Msg("saved messages not found")
		return nil, ErrSavedMessagesNotFound
	}

	logger.Info().Interface("messages", messages).Msg("successfully got saved messages")
	return messages, nil
}

func (s savedService) GetSavedMessageByMessageID(id int) (*model.Saved, error) {
	logger := s.logger

	message, err := s.store.Saved.GetSavedMessageByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("get saved message by id")
		return nil, fmt.Errorf("get saved message by id: %w", err)
	}

	if message == nil {
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
		return fmt.Errorf("create saved message: %w", err)
	}

	logger.Info().Msg("saved message successfully created")
	return nil
}

func (s savedService) DeleteSavedMessage(id int) error {
	logger := s.logger

	err := s.store.Saved.DeleteSavedMessage(id)
	if err != nil {
		logger.Error().Err(err).Msg("delete saved message")
		return fmt.Errorf("delete saved message: %w", err)
	}

	return nil
}

func (s savedService) ProcessSavedMessages(userID int) (*LoadSavedMessagesOutput, error) {
	logger := s.logger

	var savedMessages []model.FullMessage

	messages, err := s.GetSavedMessages(userID)
	if err != nil {
		if !errors.Is(err, ErrSavedMessagesNotFound) {
			logger.Error().Err(err).Msg("get saved messages")
			return nil, fmt.Errorf("get saved messages: %w", err)
		}
	}

	for _, msg := range messages {
		fullMessage, err := s.message.GetFullMessageByMessageID(msg.MessageID)
		if err != nil {
			if !errors.Is(err, ErrMessageNotFound) {
				logger.Error().Err(err).Msg("get full message by id")
			}

			continue
		}

		fullMessage.SavedID = msg.ID

		savedMessages = append(savedMessages, *fullMessage)
	}

	logger.Info().Interface("saved messages", savedMessages).Msg("successfully processed saved messages data")
	return &LoadSavedMessagesOutput{
		SavedMessages:      savedMessages,
		SavedMessagesCount: len(savedMessages),
	}, nil
}
