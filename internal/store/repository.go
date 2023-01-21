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
	GetMessagesCountByChannelID(id int) (int, error)
	GetMessageByTitle(title string) (*model.DBMessage, error)
	GetFullMessagesByPage(page int) ([]model.FullMessage, error)
	GetFullMessagesByChannelIDAndPage(id, page int) ([]model.FullMessage, error)
	GetFullMessagesByUserID(id int) ([]model.FullMessage, error)
	GetFullMessageByID(id int) (*model.FullMessage, error)
}

//go:generate mockery --dir . --name ReplyRepo --output ./mocks
type ReplyRepo interface {
	CreateReply(reply *model.DBReply) error
	GetFullRepliesByMessageID(id int) ([]model.FullReply, error)
}

//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	CreateUser(user *model.User) (int, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
}

//go:generate mockery --dir . --name WebUserRepo --output ./mocks
type WebUserRepo interface {
	GetWebUserByID(id int) (*model.WebUser, error)
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) error
}

//go:generate mockery --dir . --name SavedRepo --output ./mocks
type SavedRepo interface {
	CreateSavedMessage(saved *model.Saved) error
	GetSavedMessages(userID int) ([]model.Saved, error)
	GetSavedMessageByID(id int) (*model.Saved, error)
	DeleteSavedMessage(id int) error
}
