package service

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store"
)

type userService struct {
	store *store.Store
}

func NewUserService(store *store.Store) UserService {
	return &userService{store: store}
}

func (s userService) CreateUser(user *model.User) (int, error) {
	candidate, err := s.GetUserByUsername(user.Username)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return 0, err
		}
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

func (s userService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.store.User.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByUsername error: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s userService) GetUserByID(id int) (*model.User, error) {
	user, err := s.store.User.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByID error: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
