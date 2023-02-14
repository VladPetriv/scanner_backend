package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
	"github.com/VladPetriv/scanner_backend/pkg/password"
)

type authService struct {
	logger         *logger.Logger
	WebUserService WebUserService
}

var _ AuthService = (*authService)(nil)

func NewAuthService(webUserService WebUserService, logger *logger.Logger) *authService {
	return &authService{
		WebUserService: webUserService,
		logger:         logger,
	}
}

func (s authService) Register(user *model.WebUser) error {
	logger := s.logger

	candidate, err := s.WebUserService.GetWebUserByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, ErrWebUserNotFound) {
			logger.Error().Err(err).Msg("get web user by email")
			return fmt.Errorf("[Register] get web user by email error: %w", err)
		}
	}
	if candidate != nil {
		logger.Info().Msg("web user is exist")
		return ErrWebUserIsExist
	}

	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		logger.Error().Err(err).Msg("hash password")
		return fmt.Errorf("hash password error: %w", err)
	}

	user.Password = hashedPassword

	err = s.WebUserService.CreateWebUser(user)
	if err != nil {
		logger.Error().Err(err).Msg("create web user")
		return fmt.Errorf("[Register] create web user error: %w", err)
	}

	logger.Info().Msg("user successfully registered")
	return nil
}

func (s authService) Login(email string, userPassword string) (string, error) {
	logger := s.logger

	candidate, err := s.WebUserService.GetWebUserByEmail(email)
	if err != nil {
		if errors.Is(err, ErrWebUserNotFound) {
			logger.Error().Err(err).Msg("get web user by email")
			return "", fmt.Errorf("[Login] get web user by email error: %w", err)
		}
	}

	if password.ComparePassword(userPassword, candidate.Password) {
		logger.Info().Msg("user successfully logined")
		return email, nil
	}

	logger.Info().Msg("incorrect password")
	return "", ErrIncorrectPassword
}
