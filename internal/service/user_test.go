package service_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
)

func Test_CreateUser(t *testing.T) {
	userInput := &model.User{Username: "test", FullName: "test test", ImageURL: "test.jpg"}

	tests := []struct {
		name          string
		mock          func(userRepo *mocks.UserRepo)
		input         *model.User
		want          int
		expectedError error
	}{
		{
			name: "CreateUser successful",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, nil)
				userRepo.On("CreateUser", userInput).Return(1, nil)
			},
			input: userInput,
			want:  1,
		},
		{
			name: "CreateUser successful with founded user",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(userInput, nil)
			},
			input: userInput,
			want:  userInput.ID,
		},
		{
			name: "CreateUser failed with error when get user by username",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).
					Return(nil, fmt.Errorf("get user by username: some store error"))
			},
			input: userInput,
			expectedError: fmt.Errorf(
				"[User] Service.GetUserByUsername error: %w",
				fmt.Errorf("get user by username: some store error"),
			),
		},
		{
			name: "CreateUser failed with error when create user",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, nil)
				userRepo.On("CreateUser", userInput).Return(0, fmt.Errorf("create tg user: some store error"))
			},
			input: userInput,
			expectedError: fmt.Errorf(
				"[User] Service.CreateUser error: %w",
				fmt.Errorf("create tg user: some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.CreateUser(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
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
		name          string
		mock          func(userRepo *mocks.UserRepo)
		input         string
		want          *model.User
		expectedError error
	}{
		{
			name: "GetUserByUsername successful",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(user, nil)
			},
			input: "test",
			want:  user,
		},
		{
			name: "GetUserByUsername failed not found user",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, nil)
			},
			input:         "test",
			expectedError: service.ErrUserNotFound,
		},
		{
			name: "GetUserByUsername failed with some store error",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, fmt.Errorf("get user by username: some store error"))
			},
			input: "test",
			expectedError: fmt.Errorf(
				"[User] Service.GetUserByUsername error: %w",
				fmt.Errorf("get user by username: some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUserByUsername(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
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
		name          string
		mock          func(userRepo *mocks.UserRepo)
		input         int
		want          *model.User
		expectedError error
	}{
		{
			name: "GetUserByID successful",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(user, nil)
			},
			input: 1,
			want:  user,
		},
		{
			name: "GetUserByID failed not found user",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 404).Return(nil, nil)
			},
			input:         404,
			expectedError: service.ErrUserNotFound,
		},
		{
			name: "GetUserByID failed with some store error",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, fmt.Errorf("get user by id: some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[User] Service.GetUserByID error: %w",
				fmt.Errorf("get user by id: some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUserByID(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}
