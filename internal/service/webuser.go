package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type webUserService struct {
	store *store.Store
}

func NewWebUserService(store *store.Store) WebUserService {
	return &webUserService{store: store}
}

func (s webUserService) CreateWebUser(user *model.WebUser) error {
	err := s.store.WebUser.CreateWebUser(user)
	if err != nil {
		return fmt.Errorf("[WebUser] Service.CreateWebUser error: %w", err)
	}

	return nil
}

func (s webUserService) GetWebUserByID(userID int) (*model.WebUser, error) {
	user, err := s.store.WebUser.GetWebUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("[WebUser] Service.GetWebUserByID error: %w", err)
	}

	if user == nil {
		return nil, ErrWebUserNotFound
	}

	return user, nil
}

func (s webUserService) GetWebUserByEmail(email string) (*model.WebUser, error) {
	user, err := s.store.WebUser.GetWebUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("[WebUser] Service.GetWebUserByEmail error: %w", err)
	}

	if user == nil {
		return nil, ErrWebUserNotFound
	}

	return user, nil
}
