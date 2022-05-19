package service_test

import (
	"errors"
	"testing"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMessageService_GetMessages(t *testing.T) {
	messages := []model.Message{
		{ID: 1, ChannelID: 1, UserID: 1, Title: "test1"},
		{ID: 2, ChannelID: 1, UserID: 1, Title: "test2"},
	}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		want    []model.Message
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Message found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessages").Return(messages, nil)
			},
			want: messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessages").Return(nil, nil)
			},
			wantErr: true,
			err:     errors.New("messages not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessages").Return(nil, errors.New("error while getting messages: some error"))
			},
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessages error: error while getting messages: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessages()
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func TestMessageService_GetMessage(t *testing.T) {
	message := &model.Message{ID: 1, ChannelID: 1, UserID: 1, Title: "test"}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    *model.Message
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Message found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessage", 1).Return(message, nil)
			},
			input: 1,
			want:  message,
		},
		{
			name: "Error: [Message not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessage", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("message not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessage", 1).Return(nil, errors.New("error while getting message: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessage error: error while getting message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessage(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func TestMessageService_GetMessageByName(t *testing.T) {
	message := &model.Message{ID: 1, ChannelID: 1, UserID: 1, Title: "test"}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   string
		want    *model.Message
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Message found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").Return(message, nil)
			},
			input: "test",
			want:  message,
		},
		{
			name: "Error: [Message not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").Return(nil, nil)
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("message not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").Return(nil, errors.New("error while getting message: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessageByName error: error while getting message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessageByName(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func TestMessageService_GetFullMessages(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1"},
		{ID: 2, Title: "test2"},
	}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    []model.FullMessage
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessages", 0).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessages", 0).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("messages not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessages", 0).Return(nil, errors.New("error while getting full messages: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetFullMessages error: error while getting full messages: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessages(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func TestMessageService_GetFullMessagesByChannelID(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1", ChannelID: 1},
		{ID: 2, Title: "test2", ChannelID: 1},
	}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    []model.FullMessage
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelID", 1, 10, 0).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelID", 1, 10, 0).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("messages not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelID", 1, 10, 0).Return(nil, errors.New("error while getting full messages: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetFullMessagesByChannelID error: error while getting full messages: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessagesByChannelID(tt.input, 10, 1)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}
func TestMessageService_GetFullMessagesByUserID(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1", ChannelID: 1, UserID: 1},
		{ID: 2, Title: "test2", ChannelID: 2, UserID: 1},
	}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    []model.FullMessage
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("messages not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, errors.New("error while getting full messages by user ID: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetFullMessagesByUserID error: error while getting full messages by user ID: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessagesByUserID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}
func TestMessageService_GetFullMessagesByMessageID(t *testing.T) {
	messages := &model.FullMessage{ID: 1, Title: "test1", ChannelID: 1}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    *model.FullMessage
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByMessageID", 1).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByMessageID", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("messages not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByMessageID", 1).Return(nil, errors.New("error while getting full messages by message ID: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetFullMessageByMessageID error: error while getting full messages by message ID: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessageByMessageID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}
