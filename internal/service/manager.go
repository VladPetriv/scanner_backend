package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/store"
)

type Manager struct {
	Channel ChannelService
	Message MessageService
	Replie  ReplieService
	User    UserService
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("No store provided")
	}

	return &Manager{
		Channel: NewChannelDBService(store),
		Message: NewMessageDBService(store),
		Replie:  NewReplieDBService(store),
		User:    NewUserDBService(store),
	}, nil
}
