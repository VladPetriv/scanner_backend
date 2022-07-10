package util

import (
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
)

func CheckMessagesStatus(messages []model.FullMessage, manager *service.Manager) []model.FullMessage {
	result := make([]model.FullMessage, 0)

	for _, message := range messages {
		saved, err := manager.Saved.GetSavedMessageByMessageID(message.ID)
		if err == nil && saved != nil {
			message.Status = true
			result = append(result, message)
			continue
		}

		message.Status = false
		result = append(result, message)
	}

	return result
}
