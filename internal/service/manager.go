package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

type Manager struct {
	Channel ChannelService
	Message MessageService
	Reply   ReplyService
	User    UserService
	WebUser WebUserService
	Saved   SavedService
	Auth    AuthService
}

func NewManager(store *store.Store, logger *logger.Logger) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("no store provided")
	}

	channelService := NewChannelService(store)
	messageService := NewMessageService(store)
	replyService := NewReplyService(store)
	userService := NewUserService(store)
	webUserService := NewWebUserService(store)
	savedService := NewSavedService(store, logger, messageService)
	authService := NewAuthService(webUserService)

	srvManager := &Manager{
		Channel: channelService,
		Message: messageService,
		Reply:   replyService,
		User:    userService,
		WebUser: webUserService,
		Saved:   savedService,
		Auth:    authService,
	}

	return srvManager, nil
}
