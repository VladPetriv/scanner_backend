package util

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func CreateMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("error while hashing password")
	}

	return string(bytes), nil
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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
