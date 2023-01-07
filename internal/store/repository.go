package store

import "github.com/VladPetriv/scanner_backend/internal/model"

//go:generate mockery --dir . --name ChannelRepo --output ./mocks
type ChannelRepo interface {
	CreateChannel(channel *model.DBChannel) error
	GetChannels() ([]model.Channel, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	GetChannelStats(channelID int) (*model.Stat, error)
}

//go:generate mockery --dir . --name MessageRepo --output ./mocks
type MessageRepo interface {
	CreateMessage(message *model.DBMessage) (int, error)
	GetMessagesCount() (int, error)
	GetMessagesCountByChannelID(ID int) (int, error)
	GetFullMessagesByPage(page int) ([]model.FullMessage, error)
	GetFullMessagesByChannelIDAndPage(ID, page int) ([]model.FullMessage, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
	GetFullMessageByMessageID(ID int) (*model.FullMessage, error)
}

//go:generate mockery --dir . --name ReplyRepo --output ./mocks
type ReplyRepo interface {
	CreateReply(reply *model.DBReply) error
	GetFullRepliesByMessageID(ID int) ([]model.FullReply, error)
}

//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	CreateUser(user *model.User) (int, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(ID int) (*model.User, error)
}

//go:generate mockery --dir . --name WebUserRepo --output ./mocks
type WebUserRepo interface {
	GetWebUserByID(userID int) (*model.WebUser, error)
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
