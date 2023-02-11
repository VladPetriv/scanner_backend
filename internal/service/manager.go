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

	replyService := NewReplyService(store)
	webUserService := NewWebUserService(store)
	messageService := NewMessageService(store, logger, replyService)
	channelService := NewChannelService(store, logger, messageService)
	userService := NewUserService(store, logger, messageService)
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
