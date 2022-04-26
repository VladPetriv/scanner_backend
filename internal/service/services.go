package service

import "github.com/VladPetriv/scanner_backend/internal/model"

type ChannelService interface {
	GetChannels() ([]model.Channel, error)
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) error
	DeleteChannel(channelID int) error
}
type MessageService interface {
	GetMessages() ([]model.Message, error)
	GetMessage(messagelID int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
	CreateMessage(message *model.Message) (int, error)
	DeleteMessage(messageID int) error
}

type ReplieService interface {
	GetReplies() ([]model.Replie, error)
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
	CreateReplie(replie *model.Replie) error
	DeleteReplie(replieID int) error
}

type UserService interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) (int, error)
}
