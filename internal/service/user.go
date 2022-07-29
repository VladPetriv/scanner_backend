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

func (s *UserDBService) CreateUser(user *model.User) (int, error) {
	id, err := s.store.User.CreateUser(user)
	if err != nil {
		return id, fmt.Errorf("[User] Service.CreateUser error: %w", err)
	}

	return id, nil
}

func (s *UserDBService) GetUsers() ([]model.User, error) {
	users, err := s.store.User.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUsers error: %w", err)
	}

	if users == nil {
		return nil, fmt.Errorf("users not found")
	}

	return users, nil
}

func (s *UserDBService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.store.User.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByUsername error: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserDBService) GetUserByID(ID int) (*model.User, error) {
	user, err := s.store.User.GetUserByID(ID)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByID error: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
