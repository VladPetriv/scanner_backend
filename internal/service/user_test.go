package service_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

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
			name: "CreateUser successful with existed user",
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
					Return(nil, fmt.Errorf("some store error"))
			},
			input: userInput,
			expectedError: fmt.Errorf(
				"get user by username from db: %w", fmt.Errorf("some store error"),
			),
		},
		{
			name: "CreateUser failed with error when create user",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, nil)
				userRepo.On("CreateUser", userInput).Return(0, fmt.Errorf("some store error"))
			},
			input: userInput,
			expectedError: fmt.Errorf(
				"create user in db: %w", fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepo := &mocks.UserRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}
			log := logger.Get(&config.Config{LogLevel: "info"})

			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, log)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, log, replyService)
			userService := service.NewUserService(&store.Store{User: userRepo}, log, messageService)
			tt.mock(userRepo)

			got, err := userService.CreateUser(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			userRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_ProcessUserPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(userRepo *mocks.UserRepo, messageRepo *mocks.MessageRepo)
		input         int
		want          *service.LoadUserOutput
		expectedError error
	}{
		{
			name: "ProcessUserPage successful",
			mock: func(userRepo *mocks.UserRepo, messageRepo *mocks.MessageRepo) {
				userRepo.On("GetUserByID", 1).Return(&model.User{
					ID:       1,
					Username: "test",
					FullName: "test test",
					ImageURL: "test.jpg",
				}, nil)

				messageRepo.On("GetFullMessagesByUserID", 1).Return([]model.FullMessage{
					{ID: 1, Title: "test", UserID: 1},
					{ID: 2, Title: "test2", UserID: 1},
				}, nil)
			},
			input: 1,
			want: &service.LoadUserOutput{
				TgUser: &model.User{
					ID:       1,
					Username: "test",
					FullName: "test test",
					ImageURL: "test.jpg",
				},
				Messages: []model.FullMessage{
					{ID: 1, Title: "test", UserID: 1},
					{ID: 2, Title: "test2", UserID: 1},
				},
				MessagesCount: 2,
			},
		},
		{
			name: "ProcessUserPage failed with not found user",
			mock: func(userRepo *mocks.UserRepo, messageRepo *mocks.MessageRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrUserNotFound,
		},
		{
			name: "ProcessUserPage failed with not found user",
			mock: func(userRepo *mocks.UserRepo, messageRepo *mocks.MessageRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrUserNotFound,
		},
		{
			name: "ProcessUserPage failed with not found user",
			mock: func(userRepo *mocks.UserRepo, messageRepo *mocks.MessageRepo) {
				userRepo.On("GetUserByID", 1).Return(&model.User{ID: 1}, nil)
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, nil)
			},
			input: 1,
			want: &service.LoadUserOutput{
				TgUser: &model.User{ID: 1},
			},
		},
		{
			name: "ProcessUserPage failed with some store error while get user by id",
			mock: func(userRepo *mocks.UserRepo, messageRepo *mocks.MessageRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get user by id from db: %w", fmt.Errorf("some store error")),
		},
		{
			name: "ProcessUserPage failed with some store error while get messages user by id",
			mock: func(userRepo *mocks.UserRepo, messageRepo *mocks.MessageRepo) {
				userRepo.On("GetUserByID", 1).Return(&model.User{ID: 1}, nil)
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[ProcessUserPage]: %w",
				fmt.Errorf("get full messages by user id from db: %w",
					fmt.Errorf("some store error"),
				),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepo := &mocks.UserRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}
			log := logger.Get(&config.Config{LogLevel: "info"})

			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, log)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, log, replyService)
			userService := service.NewUserService(&store.Store{User: userRepo}, log, messageService)
			tt.mock(userRepo, messageRepo)

			got, err := userService.ProcessUserPage(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			userRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}
