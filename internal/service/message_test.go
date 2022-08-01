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

func Test_CreateMesage(t *testing.T) {
	messageInput := &model.DBMessage{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.url", ImageURL: "test.jpg"}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   *model.DBMessage
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "OK: [Message created]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("CreateMessage", messageInput).Return(1, nil)
			},
			input: messageInput,
			want:  1,
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("CreateMessage", messageInput).Return(0, errors.New("failed to create message: some error"))
			},
			input:   messageInput,
			wantErr: true,
			err:     errors.New("[Message] Service.CreateMessage error: failed to create message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.CreateMessage(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.EqualValues(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func Test_GetMessagesCount(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages count found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(10, nil)
			},
			want: 10,
		},
		{
			name: "Error: [Messages count not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, pg.ErrMessagesCountNotFound)
			},
			wantErr: true,
			err:     fmt.Errorf("[Message] Service.GetMessagesCount error: %w", pg.ErrMessagesCountNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, errors.New("error while getting messages count: some error"))
			},
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessagesCount error: error while getting messages count: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessagesCount()
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

func Test_GetMessagesLengthByChannelID(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages count found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(3, nil)
			},
			input: 1,
			want:  3,
		},
		{
			name: "Error: [Messages count not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(0, pg.ErrMessagesCountNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[Message] Service.GetMessagesCountByChannelID error: %w", pg.ErrMessagesCountNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(0, errors.New("error while getting  messages count by channel ID: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessagesCountByChannelID error: error while getting  messages count by channel ID: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessagesCountByChannelID(tt.input)
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

func Test_GetMessagesCountByChannelID(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages count found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(10, nil)
			},
			input: 1,
			want:  10,
		},
		{
			name: "Error: [Messages count not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(0, pg.ErrMessagesCountNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[Message] Service.GetMessagesCountByChannelID error: %w", pg.ErrMessagesCountNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(0, errors.New("error while getting messages count by channel id: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessagesCountByChannelID error: error while getting messages count by channel id: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessagesCountByChannelID(tt.input)
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

func Test_GetFullMessagesByPage(t *testing.T) {
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
				messageRepo.On("GetFullMessagesByPage", 0).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 0).Return(nil, pg.ErrMessagesNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[Message] Service.GetFullMessagesByPage error: %w", pg.ErrMessagesNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 0).Return(nil, errors.New("error while getting full messages by page: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetFullMessagesByPage error: error while getting full messages by page: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessagesByPage(tt.input)
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

func Test_GetFullMessagesByChannelIDAndPage(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1", ChannelID: 1, ImageURL: "test1.jpg"},
		{ID: 2, Title: "test2", ChannelID: 1, ImageURL: "test2.jpg"},
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
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(messages, nil)
			},
			input: 1,
			want:  messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(nil, pg.ErrMessagesNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[Message] Service.GetFullMessagesByChannelIDAndPage error: %w", pg.ErrMessagesNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(nil, errors.New("error while getting full messages by channel ID: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetFullMessagesByChannelIDAndPage error: error while getting full messages by channel ID: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetFullMessagesByChannelIDAndPage(tt.input, 1)
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

func Test_GetFullMessagesByUserID(t *testing.T) {
	messages := []model.FullMessage{
		{ID: 1, Title: "test1", ChannelID: 1, UserID: 1, ImageURL: "test1.jpg"},
		{ID: 2, Title: "test2", ChannelID: 2, UserID: 1, ImageURL: "test2.jpg"},
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
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, pg.ErrMessagesNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[Message] Service.GetFullMessagesByUserID error: %w", pg.ErrMessagesNotFound),
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

func Test_GetFullMessagesByMessageID(t *testing.T) {
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
				messageRepo.On("GetFullMessageByMessageID", 1).Return(nil, pg.ErrMessageNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[Message] Service.GetFullMessageByMessageID error: %w", pg.ErrMessageNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByMessageID", 1).Return(nil, errors.New("error while getting full message by message ID: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetFullMessageByMessageID error: error while getting full message by message ID: some error"),
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
