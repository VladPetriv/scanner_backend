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
	WebUser WebUserService
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("no store provided")
	}

	return &Manager{
		Channel: NewChannelDBService(store),
		Message: NewMessageDBService(store),
		Replie:  NewReplieDBService(store),
		User:    NewUserDBService(store),
		WebUser: NewWebUserDbService(store),
	}, nil
}
