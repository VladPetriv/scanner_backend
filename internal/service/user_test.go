package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
)

func Test_CreateUser(t *testing.T) {
	userInput := &model.User{Username: "test", FullName: "test test", ImageURL: "test.jpg"}
	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		input   *model.User
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User created]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, pg.ErrUserNotFound)
				userRepo.On("CreateUser", userInput).Return(1, nil)
			},
			input: userInput,
			want:  1,
		},
		{
			name: "Error: [User is exist]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(userInput, nil)
			},
			input: userInput,
			want:  userInput.ID,
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, pg.ErrUserNotFound)
				userRepo.On("CreateUser", userInput).Return(0, errors.New("failed to create user: some error"))
			},
			input:   userInput,
			wantErr: true,
			err:     errors.New("[User] Service.CreateUser error: failed to create user: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.CreateUser(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.EqualValues(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}

func Test_GetUserByUsername(t *testing.T) {
	user := &model.User{ID: 1, Username: "test1", FullName: "test1 test", ImageURL: "test1.jpg"}

	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		input   string
		want    *model.User
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(user, nil)
			},
			input: "test",
			want:  user,
		},
		{
			name: "Error: [User not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, pg.ErrUserNotFound)
			},
			input:   "test",
			wantErr: true,
			err:     fmt.Errorf("[User] Service.GetUserByUsername error: %w", pg.ErrUserNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, errors.New("error while getting users: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[User] Service.GetUserByUsername error: error while getting users: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUserByUsername(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}

func Test_GetUserByID(t *testing.T) {
	user := &model.User{ID: 1, Username: "test1", FullName: "test1 test", ImageURL: "test1.jpg"}

	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		input   int
		want    *model.User
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(user, nil)
			},
			input: 1,
			want:  user,
		},
		{
			name: "Error: [User not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, pg.ErrUserNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[User] Service.GetUserByID error: %w", pg.ErrUserNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, errors.New("error while getting users: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[User] Service.GetUserByID error: error while getting users: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUserByID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}
