package store

import "github.com/VladPetriv/scanner_backend/internal/model"

type ChannelRepo interface {
	GetChannels() ([]model.Channel, error)
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}

type MessageRepo interface {
	GetMessages() ([]model.Message, error)
	GetMessage(messageID int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
	GetFullMessages() ([]model.FullMessage, error)
	GetFullMessagesByChannelID(ID int) ([]model.FullMessage, error)
}

type ReplieRepo interface {
	GetReplies() ([]model.Replie, error)
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
}

type UserRepo interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
}
