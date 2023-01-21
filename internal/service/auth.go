package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/password"
)

type authService struct {
	WebUserService WebUserService
}

func NewAuthService(webUserService WebUserService) AuthService {
	return &authService{
		WebUserService: webUserService,
	}
}

func (s authService) Register(user *model.WebUser) error {
	candidate, err := s.WebUserService.GetWebUserByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, ErrWebUserNotFound) {
			return fmt.Errorf("[Auth] Service.Register error: %w", err)
		}
	}

	if candidate != nil {
		return ErrWebUserIsExist
	}

	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("[Auth] Service.Register error: %w", err)
	}

	user.Password = hashedPassword

	err = s.WebUserService.CreateWebUser(user)
	if err != nil {
		return fmt.Errorf("[Auth] Service.Register error: %w", err)
	}

	return nil
}

func (s authService) Login(email string, userPassword string) (string, error) {
	candidate, err := s.WebUserService.GetWebUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("[Auth] Service.Login error: %w", err)
	}

	if password.ComparePassword(userPassword, candidate.Password) {
		return email, nil
	}

	return "", ErrIncorrectPassword
}
