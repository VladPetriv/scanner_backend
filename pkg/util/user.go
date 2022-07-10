package util

import "github.com/VladPetriv/scanner_backend/internal/model"

func ProcessWebUserData(user *model.WebUser) (int, string) {
	if user != nil {
		return user.ID, user.Email
	}

	return 0, ""
}
