package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
)

type UserDBService struct {
	store *store.Store
}

func NewUserDBService(store *store.Store) *UserDBService {
	return &UserDBService{store: store}
}

func (s *UserDBService) CreateUser(user *model.User) (int, error) {
	candidate, err := s.GetUserByUsername(user.Username)
	if !errors.Is(err, pg.ErrUserNotFound) {
		return 0, err
	}

	if candidate != nil {
		return candidate.ID, nil
	}

	id, err := s.store.User.CreateUser(user)
	if err != nil {
		return id, fmt.Errorf("[User] Service.CreateUser error: %w", err)
	}

	return id, nil
}

func (s *UserDBService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.store.User.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByUsername error: %w", err)
	}

	return user, nil
}

func (s *UserDBService) GetUserByID(ID int) (*model.User, error) {
	user, err := s.store.User.GetUserByID(ID)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByID error: %w", err)
	}

	return user, nil
}
