package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type UserDBService struct {
	store *store.Store
}

func NewUserDBService(store *store.Store) *UserDBService {
	return &UserDBService{store: store}
}

func (s *UserDBService) GetUsers() ([]model.User, error) {
	users, err := s.store.User.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUser error: %w", err)
	}

	if users == nil {
		return nil, fmt.Errorf("user not found")
	}

	return users, nil
}

func (s *UserDBService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.store.User.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByUsername error: %w", err)
	}
	if user == nil {
		return nil, nil
	}

	return user, nil
}
