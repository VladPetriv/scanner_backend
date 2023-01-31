package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type userService struct {
	logger  *logger.Logger
	message MessageService
	store   *store.Store
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

	candidate, err := s.GetUserByUsername(user.Username)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			logger.Error().Err(err).Msg("get user by username")
			return 0, err
		}
	}
	if candidate != nil {
		return candidate.ID, nil
	}

	id, err := s.store.User.CreateUser(user)
	if err != nil {
		logger.Error().Err(err).Msg("create user")
		return id, fmt.Errorf("create user: %w", err)
	}

	logger.Info().Int("userID", id).Msg("user successfully created")
	return id, nil
}

func (s userService) GetUserByUsername(username string) (*model.User, error) {
	logger := s.logger

	user, err := s.store.User.GetUserByUsername(username)
	if err != nil {
		logger.Error().Err(err).Msg("get user by username")
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	logger.Info().Interface("user", user).Msg("successfully got user by username")
	return user, nil
}

func (s userService) GetUserByID(id int) (*model.User, error) {
	logger := s.logger

	user, err := s.store.User.GetUserByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("get user by id")
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	logger.Info().Interface("user", user).Msg("successfully got user by id")
	return user, nil
}

func (s userService) ProcessUserPage(userID int) (*LoadUserOutput, error) {
	logger := s.logger

	user, err := s.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			logger.Info().Err(err).Msg("user not found")
			return nil, err
		}

		logger.Error().Err(err).Msg("get user by id")
		return nil, err
	}

	messages, err := s.message.GetFullMessagesByUserID(user.ID)
	if err != nil {
		if errors.Is(err, ErrMessagesNotFound) {
			return nil, err
		}

		logger.Error().Err(err).Msg("get messages by user id")
		return nil, fmt.Errorf("get messages by user id: %w", err)
	}

	logger.Info().Msg("successfully processed user page")
	return &LoadUserOutput{
		TgUser:        user,
		Messages:      messages,
		MessagesCount: len(messages),
	}, nil
}
