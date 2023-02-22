package service

import (
	"errors"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type ChannelService interface {
	CreateChannel(channel *model.DBChannel) error
	GetChannels() ([]model.Channel, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	GetChannelStats(channelID int) (*model.Stat, error)
	ProcessChannelPage(channelName string, page int) (*LoadChannelOutput, error)
	ProcessChannelsPage(page int) (*LoadChannelsOutput, error)
}

type LoadChannelOutput struct {
	Channel       model.Channel
	Messages      []model.FullMessage
	MessagesCount int
}

type LoadChannelsOutput struct {
	Channels []model.Channel
}

var (
	ErrChannelsNotFound         = errors.New("channels not found")
	ErrChannelNotFound          = errors.New("channel not found")
	ErrChannelExists            = errors.New("channel is exist")
	ErrChannelStatisticNotFound = errors.New("channel statistic not found")
)

type MessageService interface {
	CreateMessage(message *model.DBMessage) (int, error)
	GetMessagesCount() (int, error)
	GetMessagesCountByChannelID(ID int) (int, error)
	GetFullMessagesByChannelIDAndPage(ID, page int) ([]model.FullMessage, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
	GetFullMessageByMessageID(id int) (*model.FullMessage, error)
	ProcessMessagePage(messageID int) (*LoadMessageOutput, error)
	ProcessHomePage(page int) (*LoadHomeOutput, error)
}

type LoadMessageOutput struct {
	Message *model.FullMessage
}

type LoadHomeOutput struct {
	Messages      []model.FullMessage
	MessagesCount int
}

var (
	ErrMessagesCountNotFound = errors.New("message count not found")
	ErrMessagesNotFound      = errors.New("messages not found")
	ErrMessageNotFound       = errors.New("messages not found")
	ErrMessageExists         = errors.New("message is exist")
)

type ReplyService interface {
	CreateReply(reply *model.DBReply) error
	GetFullRepliesByMessageID(ID int) ([]model.FullReply, error)
}

var ErrRepliesNotFound = errors.New("replies not found")

type SavedService interface {
	GetSavedMessageByMessageID(id int) (*model.Saved, error)
	CreateSavedMessage(saved *model.Saved) error
	DeleteSavedMessage(ID int) error
	ProcessSavedMessages(userID int) (*LoadSavedMessagesOutput, error)
}

type LoadSavedMessagesOutput struct {
	SavedMessages      []model.FullMessage
	SavedMessagesCount int
}

var (
	ErrSavedMessagesNotFound = errors.New("saved messages not found")
	ErrSavedMessageNotFound  = errors.New("saved message not found")
)

type UserService interface {
	CreateUser(user *model.User) (int, error)
	ProcessUserPage(userID int) (*LoadUserOutput, error)
}

type LoadUserOutput struct {
	TgUser        *model.User
	Messages      []model.FullMessage
	MessagesCount int
}

var ErrUserNotFound = errors.New("user not found")

type WebUserService interface {
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) error
}

var ErrWebUserNotFound = errors.New("web user not found")

type AuthService interface {
	Login(email string, userPassword string) (string, error)
	Register(user *model.WebUser) error
}

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrWebUserIsExist    = errors.New("web user is exist")
)
