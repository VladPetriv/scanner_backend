package util

import "github.com/VladPetriv/scanner_backend/internal/model"

func ProcessChannels(channels []model.Channel) []model.Channel {
	if len(channels) <= 10 {
		return channels
	} else {
		return channels[:10]
	}
}
