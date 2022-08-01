package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type WebUserDbService struct {
	store *store.Store
}

func NewWebUserDbService(store *store.Store) *WebUserDbService {
	return &WebUserDbService{store: store}
}

func (s *WebUserDbService) GetWebUserByID(userID int) (*model.WebUser, error) {
	user, err := s.store.WebUser.GetWebUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("[WebUser] Service.GetWebUserByID error: %w", err)
	}

	return user, nil
}

func (s *WebUserDbService) GetWebUserByEmail(email string) (*model.WebUser, error) {
	user, err := s.store.WebUser.GetWebUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("[WebUser] Service.GetWebUserByEmail error: %w", err)
	}

	return user, nil
}

func (s *WebUserDbService) CreateWebUser(user *model.WebUser) error {
	_, err := s.store.WebUser.CreateWebUser(user)
	if err != nil {
		return fmt.Errorf("[WebUser] Service.CreateWebUser error: %w", err)
	}

	return nil
}
