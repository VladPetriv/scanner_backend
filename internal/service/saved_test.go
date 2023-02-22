package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestSavedService_GetSavedMessageByMessageID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(savedRepo *mocks.SavedRepo)
		input         int
		want          *model.Saved
		expectedError error
	}{
		{
			name: "GetSavedMessageByMessageID successful",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByID", 1).Return(&model.Saved{ID: 1, WebUserID: 1, MessageID: 1}, nil)
			},
			input: 1,
			want:  &model.Saved{ID: 1, WebUserID: 1, MessageID: 1},
		},
		{
			name: "GetSavedMessageByMessageID failed with not found message",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrSavedMessageNotFound,
		},
		{
			name: "GetSavedMessageByMessageID failed with some store erorr",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByID", 1).Return(nil, errors.New("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get saved message by message id from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			savedRepo := &mocks.SavedRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			savedService := service.NewSavedService(&store.Store{Saved: savedRepo}, logger, messageService)
			tt.mock(savedRepo)

			got, err := savedService.GetSavedMessageByMessageID(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			savedRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestSavedService_CreateSavedMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(savedRepo *mocks.SavedRepo)
		input         *model.Saved
		expectedError error
	}{
		{
			name: "CreateSavedMessage successful",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", &model.Saved{WebUserID: 1, MessageID: 1}).Return(nil)
			},
			input: &model.Saved{WebUserID: 1, MessageID: 1},
		},
		{
			name: "CreateSavedMessage failed with some store error",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", &model.Saved{
					WebUserID: 1,
					MessageID: 1,
				}).Return(fmt.Errorf("some store error"))
			},
			input: &model.Saved{WebUserID: 1, MessageID: 1},
			expectedError: fmt.Errorf(
				"create saved message in db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			savedRepo := &mocks.SavedRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			savedService := service.NewSavedService(&store.Store{Saved: savedRepo}, logger, messageService)
			tt.mock(savedRepo)

			err := savedService.CreateSavedMessage(tt.input)
			assert.Equal(t, tt.expectedError, err)

			savedRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestSavedService_DeleteSavedMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(savedRepo *mocks.SavedRepo)
		input         int
		expectedError error
	}{
		{
			name: "DeleteSavedMessage successful",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("DeleteSavedMessage", 1).Return(nil)
			},
			input: 1,
		},
		{
			name: "DeleteSavedMessage failed with some store error",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("DeleteSavedMessage", 1).Return(fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"delete saved message from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			savedRepo := &mocks.SavedRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			savedService := service.NewSavedService(&store.Store{Saved: savedRepo}, logger, messageService)
			tt.mock(savedRepo)

			err := savedService.DeleteSavedMessage(tt.input)
			assert.EqualValues(t, tt.expectedError, err)

			savedRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestSavedService_ProcessSavedMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         int
		mock          func(savedRepo *mocks.SavedRepo, messageRepo *mocks.MessageRepo)
		want          *service.LoadSavedMessagesOutput
		expectedError error
	}{
		{
			name:  "ProcessSavedMessages successful",
			input: 1,
			mock: func(savedRepo *mocks.SavedRepo, messageRepo *mocks.MessageRepo) {
				savedRepo.On("GetSavedMessages", 1).Return([]model.Saved{
					{MessageID: 1, WebUserID: 1},
					{MessageID: 2, WebUserID: 1},
				}, nil)

				messageRepo.On("GetFullMessageByID", 1).Return(&model.FullMessage{
					ID: 1,
				}, nil)

				messageRepo.On("GetFullMessageByID", 2).Return(&model.FullMessage{
					ID: 2,
				}, nil)
			},
			want: &service.LoadSavedMessagesOutput{
				SavedMessages: []model.FullMessage{
					{ID: 1}, {ID: 2},
				},
				SavedMessagesCount: 2,
			},
		},
		{
			name:  "ProcessSavedMessages failed with not found saved messages",
			input: 1,
			mock: func(savedRepo *mocks.SavedRepo, messageRepo *mocks.MessageRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, nil)
			},
			want: &service.LoadSavedMessagesOutput{},
		},
		{
			name:  "ProcessSavedMessages failed with some store error when get saved messages",
			input: 1,
			mock: func(savedRepo *mocks.SavedRepo, messageRepo *mocks.MessageRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, fmt.Errorf("some store error"))
			},
			expectedError: fmt.Errorf(
				"get saved messages from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			savedRepo := &mocks.SavedRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			savedService := service.NewSavedService(&store.Store{Saved: savedRepo}, logger, messageService)
			tt.mock(savedRepo, messageRepo)

			got, err := savedService.ProcessSavedMessages(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			savedRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}
