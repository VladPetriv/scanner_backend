package store

import "github.com/VladPetriv/scanner_backend/internal/model"

//go:generate mockery --dir . --name ChannelRepo --output ./mocks
type ChannelRepo interface {
	GetChannels() ([]model.Channel, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	GetChannelStats(channelID int) (*model.Stat, error)
}

//go:generate mockery --dir . --name MessageRepo --output ./mocks
type MessageRepo interface {
	GetMessagesLength() (int, error)
	GetFullMessages(page int) ([]model.FullMessage, error)
	GetFullMessagesByChannelID(ID, limit, page int) ([]model.FullMessage, error)
	GetMessagesLengthByChannelID(ID int) (int, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
	GetFullMessageByMessageID(ID int) (*model.FullMessage, error)
}

//go:generate mockery --dir . --name ReplieRepo --output ./mocks
type ReplieRepo interface {
	GetReplies() ([]model.Replie, error)
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
	GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error)
}

//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(ID int) (*model.User, error)
}

//go:generate mockery --dir . --name WebUserRepo --output ./mocks
type WebUserRepo interface {
	GetWebUser(userID int) (*model.WebUser, error)
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) (int, error)
}

//go:generate mockery --dir . --name SavedRepo --output ./mocks
type SavedRepo interface {
	GetSavedMessages(UserID int) ([]model.Saved, error)
	GetSavedMessageByMessageID(ID int) (*model.Saved, error)
	CreateSavedMessage(saved *model.Saved) (int, error)
	DeleteSavedMessage(ID int) (int, error)
}
