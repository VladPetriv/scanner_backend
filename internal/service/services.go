package service

import "github.com/VladPetriv/scanner_backend/internal/model"

type ChannelService interface {
	CreateChannel(channel *model.DBChannel) error
	GetChannels() ([]model.Channel, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	GetChannelStats(channelID int) (*model.Stat, error)
}

type MessageService interface {
	CreateMessage(message *model.DBMessage) (int, error)
	GetMessagesLength() (int, error)
	GetFullMessages(page int) ([]model.FullMessage, error)
	GetFullMessagesByChannelID(ID, limit, page int) ([]model.FullMessage, error)
	GetMessagesLengthByChannelID(ID int) (int, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
	GetFullMessageByMessageID(ID int) (*model.FullMessage, error)
}

type ReplieService interface {
	CreateReplie(replie *model.DBReplie) error
	GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error)
}

type UserService interface {
	CreateUser(user *model.User) (int, error)
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(ID int) (*model.User, error)
}

type WebUserService interface {
	GetWebUser(userID int) (*model.WebUser, error)
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) error
}

type SavedService interface {
	GetSavedMessages(UserID int) ([]model.Saved, error)
	GetSavedMessageByMessageID(ID int) (*model.Saved, error)
	CreateSavedMessage(saved *model.Saved) error
	DeleteSavedMessage(ID int) error
}
