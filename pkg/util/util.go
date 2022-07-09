package util

import (
	"github.com/VladPetriv/scanner_backend/internal/model"
)

func ProcessChannels(channels []model.Channel) []model.Channel {
	if len(channels) <= 10 {
		return channels
	} else {
		return channels[:10]
	}
}

func ProcessWebUserData(user *model.WebUser) (int, string) {
	if user != nil {
		return user.ID, user.Email
	}

	return 0, ""
}
