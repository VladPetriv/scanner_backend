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

func Test_CreateMessage(t *testing.T) {
	messageInput := &model.DBMessage{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.url", ImageURL: "test.jpg"}

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         *model.DBMessage
		want          int
		expectedError error
	}{
		{
			name: "CreateMessage successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(nil, nil)
				messageRepo.On("CreateMessage", messageInput).Return(1, nil)
			},
			input: messageInput,
			want:  1,
		},
		{
			name: "CreateMessage failed with existed message",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(messageInput, nil)
			},
			input:         messageInput,
			expectedError: service.ErrMessageExists,
		},
		{
			name: "CreateMessage failed with some store error when get message by title",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(nil, fmt.Errorf("get message by title: some error"))
			},
			input: messageInput,
			expectedError: fmt.Errorf(
				"[Message] Service.GetMessageByTitle error: %w",
				fmt.Errorf("get message by title: some error"),
			),
		},
		{
			name: "CreateMessage failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(nil, nil)
				messageRepo.On("CreateMessage", messageInput).Return(0, fmt.Errorf("create message: some error"))
			},
			input: messageInput,
			expectedError: fmt.Errorf(
				"[Message] Service.CreateMessage error: %w",
				fmt.Errorf("create message: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.CreateMessage(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetMessagesCount(t *testing.T) {
	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		want          int
		expectedError error
	}{
		{
			name: "GetMessageCount successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(10, nil)
			},
			want: 10,
		},
		{
			name: "GetMessagesCount failed with not found messages count",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, nil)
			},
			expectedError: service.ErrMessagesCountNotFound,
		},
		{
			name: "GetMessagesCount failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, fmt.Errorf("get messages count: some error"))
			},
			expectedError: fmt.Errorf(
				"[Message] Service.GetMessagesCount error: %w",
				fmt.Errorf("get messages count: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessagesCount()
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetMessagesLengthByChannelID(t *testing.T) {
	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         int
		want          int
		expectedError error
	}{
		{
			name: "GetMessagesCountByChannelID successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(3, nil)
			},
			input: 1,
			want:  3,
		},
		{
			name: "GetMessagesCountByChannelID failed with not found messages count",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(0, nil)
			},
			input:         1,
			expectedError: service.ErrMessagesCountNotFound,
		},
		{
			name: "GetMessagesCountByChannelID failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).
					Return(0, fmt.Errorf("get messages count by channel id: some error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Message] Service.GetMessagesCountByChannelID error: %w",
				fmt.Errorf("get messages count by channel id: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessagesCountByChannelID(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetMessageByTitle(t *testing.T) {
	message := &model.DBMessage{ID: 1, Title: "test"}

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         string
		want          *model.DBMessage
		expectedError error
	}{
		{
			name: "GetMessageByTitle successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(message, nil)
			},
			input: "test",
			want:  message,
		},
		{
			name: "GetMessageByTitle failed with not found message",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(nil, nil)
			},
			input:         "test",
			expectedError: service.ErrMessageNotFound,
		},
		{
			name: "GetMessageByTitle failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(nil, fmt.Errorf("get message by title: some error"))
			},
			input: "test",
			expectedError: fmt.Errorf(
				"[Message] Service.GetMessageByTitle error: %w",
				fmt.Errorf("get message by title: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessageByTitle(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetFullMessagesByPage(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1"},
		{ID: 2, Title: "test2"},
	}

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         int
		want          []model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessagesByPage successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 0).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "GetFullMessagesByPage failed with not found messages",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 0).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrMessagesNotFound,
		},
		{
			name: "GetFullMessagesByPage failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 0).Return(nil, fmt.Errorf("get full messages by page: some error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Message] Service.GetFullMessagesByPage error: %w",
				fmt.Errorf("get full messages by page: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessagesByPage(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetFullMessagesByChannelIDAndPage(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1", ChannelID: 1, ImageURL: "test1.jpg"},
		{ID: 2, Title: "test2", ChannelID: 1, ImageURL: "test2.jpg"},
	}

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         int
		want          []model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessagesByChannelIDAndPage successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "GetFullMessagesByChannelIDAndPage failed with not found messages",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrMessagesNotFound,
		},
		{
			name: "GetFullMessagesByChannelIDAndPage failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).
					Return(nil, fmt.Errorf("get full messages by channel id: some error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Message] Service.GetFullMessagesByChannelIDAndPage error: %w",
				fmt.Errorf("get full messages by channel id: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessagesByChannelIDAndPage(tt.input, 1)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetFullMessagesByUserID(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1", ChannelID: 1, UserID: 1, ImageURL: "test1.jpg"},
		{ID: 2, Title: "test2", ChannelID: 2, UserID: 1, ImageURL: "test2.jpg"},
	}

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         int
		want          []model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessagesByUserID successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "GetFullMessagesByUserID failed with not found messages",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrMessagesNotFound,
		},
		{
			name: "GetFullMessagesByUserID failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, fmt.Errorf("get full messages by user id: some error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Message] Service.GetFullMessagesByUserID error: %w",
				fmt.Errorf("get full messages by user id: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessagesByUserID(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetFullMessagesByMessageID(t *testing.T) {
	messages := &model.FullMessage{ID: 1, Title: "test1", ChannelID: 1}

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         int
		want          *model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessagesByMessageID successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "GetFullMessagesByMessageID failed with not found message",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrMessageNotFound,
		},
		{
			name: "GetFullMessagesByMessageID failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(nil, fmt.Errorf("get full message by message id: some error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Message] Service.GetFullMessageByMessageID error: %w",
				fmt.Errorf("get full message by message id: some error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessageByMessageID(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}
