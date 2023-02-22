package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type webUserService struct {
	store  *store.Store
	logger *logger.Logger
}

var _ WebUserService = (*webUserService)(nil)

func NewWebUserService(store *store.Store, logger *logger.Logger) *webUserService {
	return &webUserService{
		store:  store,
		logger: logger,
	}
}

func (s webUserService) CreateWebUser(user *model.WebUser) error {
	logger := s.logger

	err := s.store.WebUser.CreateWebUser(user)
	if err != nil {
		logger.Error().Err(err).Msg("create web user")
		return fmt.Errorf("create web user in db: %w", err)
	}

	logger.Info().Msg("web user successfully created")
	return nil
}

func (s webUserService) GetWebUserByEmail(email string) (*model.WebUser, error) {
	logger := s.logger

	user, err := s.store.WebUser.GetWebUserByEmail(email)
	if err != nil {
		logger.Error().Err(err).Msg("get web user by email")
		return nil, fmt.Errorf("get web user by email from db: %w", err)
	}
	if user == nil {
		logger.Info().Str("user email", email).Msg("web user by email not found")
		return nil, ErrWebUserNotFound
	}

	logger.Info().Interface("web user", user).Msg("successfully got web user")
	return user, nil
}
