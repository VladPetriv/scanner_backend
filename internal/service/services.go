package service

import "github.com/VladPetriv/scanner_backend/internal/model"

type ChannelService interface {
	GetChannels() ([]model.Channel, error)
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}
type MessageService interface {
	GetMessages() ([]model.Message, error)
	GetMessage(messagelID int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
}

type ReplieService interface {
	GetReplies() ([]model.Replie, error)
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
}

type UserService interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
}
