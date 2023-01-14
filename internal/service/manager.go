package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/store"
)

type Manager struct {
	Channel ChannelService
	Message MessageService
	Reply   ReplyService
	User    UserService
	WebUser WebUserService
	Saved   SavedService
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("no store provided")
	}

	return &Manager{
		Channel: NewChannelService(store),
		Message: NewMessageService(store),
		Reply:   NewReplyService(store),
		User:    NewUserService(store),
		WebUser: NewWebUserService(store),
		Saved:   NewSavedService(store),
	}, nil
}
