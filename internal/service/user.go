package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type userService struct {
	store   *store.Store
	logger  *logger.Logger
	message MessageService
}

var _ UserService = (*userService)(nil)

func NewUserService(store *store.Store, logger *logger.Logger, message MessageService) *userService {
	return &userService{
		store:   store,
		logger:  logger,
		message: message,
	}
}

func (s userService) CreateUser(user *model.User) (int, error) {
	logger := s.logger

	candidate, err := s.store.User.GetUserByUsername(user.Username)
	if err != nil {
		logger.Error().Err(err).Msg("get user by username")
		return 0, fmt.Errorf("get user by username from db: %w", err)
	}
	if candidate != nil {
		logger.Info().Msg("user found, don't create new user")
		return candidate.ID, nil
	}

	id, err := s.store.User.CreateUser(user)
	if err != nil {
		logger.Error().Err(err).Msg("create user")
		return id, fmt.Errorf("create user in db: %w", err)
	}

	logger.Info().Int("userID", id).Msg("user successfully created")
	return id, nil
}

func (s userService) ProcessUserPage(userID int) (*LoadUserOutput, error) {
	logger := s.logger

	user, err := s.store.User.GetUserByID(userID)
	if err != nil {
		logger.Error().Err(err).Msg("get user by id")
		return nil, fmt.Errorf("get user by id from db: %w", err)
	}
	if user == nil {
		logger.Info().Msg("user not found")
		return nil, ErrUserNotFound
	}

	messages, err := s.message.GetFullMessagesByUserID(user.ID)
	if err != nil {
		if errors.Is(err, ErrMessagesNotFound) {
			logger.Info().Int("user id", user.ID).Msg("messages by user id not found")
			return &LoadUserOutput{
				TgUser: user,
			}, nil
		}

		logger.Error().Err(err).Msg("get messages by user id")
		return nil, fmt.Errorf("[ProcessUserPage]: %w", err)
	}

	return &LoadUserOutput{
		TgUser:        user,
		Messages:      messages,
		MessagesCount: len(messages),
	}, nil
}
