package service

import "github.com/VladPetriv/scanner_backend/internal/model"

type ChannelService interface {
	GetChannels() ([]model.Channel, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}
type MessageService interface {
	GetMessages() ([]model.Message, error)
	GetMessage(messagelID int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
	GetFullMessages(page int) ([]model.FullMessage, error)
	GetFullMessagesByChannelID(ID, limit, page int) ([]model.FullMessage, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
	GetFullMessageByMessageID(ID int) (*model.FullMessage, error)
}

type ReplieService interface {
	GetReplies() ([]model.Replie, error)
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
	GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error)
}

type UserService interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(ID int) (*model.User, error)
}

type WebUserService interface {
	GetWebUser(userID int) (*model.WebUser, error)
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) error
}
