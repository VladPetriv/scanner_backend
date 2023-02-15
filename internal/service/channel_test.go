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

func Test_CreateChannel(t *testing.T) {
	channelInput := &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"}

	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		input         *model.DBChannel
		expectedError error
	}{
		{
			name: "CreateChannel successful",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", channelInput.Name).Return(nil, nil)
				channelRepo.On("CreateChannel", channelInput).Return(nil)
			},
			input: channelInput,
		},
		{
			name: "CreateChannel failed with existed channel",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", channelInput.Name).Return(&model.Channel{Name: "test"}, nil)
			},
			input:         channelInput,
			expectedError: service.ErrChannelExists,
		},
		{
			name: "CreateChannel failed with some store error when get channel by name",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", channelInput.Name).
					Return(nil, fmt.Errorf("get channel by name: some store error"))
			},
			input: channelInput,
			expectedError: fmt.Errorf(
				"[Channel] Service.GetChannelByName error: %w",
				fmt.Errorf("get channel by name: some store error"),
			),
		},
		{
			name: "CreateChannel failed with some store error when create channel",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", channelInput.Name).Return(nil, nil)
				channelRepo.On("CreateChannel", channelInput).Return(fmt.Errorf("create channel: some store error"))
			},
			input: channelInput,
			expectedError: fmt.Errorf(
				"[Channel] Service.CreateChannel error: %w",
				fmt.Errorf("create channel: some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		messageRepo := &mocks.MessageRepo{}
		replyRepo := &mocks.ReplyRepo{}

		logger := logger.Get(&config.Config{LogLevel: "info"})
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
		messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
		channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
		tt.mock(channelRepo)

		err := channelService.CreateChannel(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
		}

		channelRepo.AssertExpectations(t)
	}
}

func Test_GetChannels(t *testing.T) {
	data := []model.Channel{
		{ID: 1, Name: "test1", Title: "test1", ImageURL: "test1.jpg"},
		{ID: 2, Name: "test2", Title: "test2", ImageURL: "test2.jpg"},
	}

	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		want          []model.Channel
		expectedError error
	}{
		{
			name: "GetChannels successful",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannels").Return(data, nil)
			},
			want: data,
		},
		{
			name: "GetChannel failed with not found channels",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannels").Return(nil, nil)
			},
			expectedError: service.ErrChannelsNotFound,
		},
		{
			name: "GetChannel failed with some store error",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannels").Return(nil, fmt.Errorf("get channels: some error"))
			},
			expectedError: fmt.Errorf("[Channel] Service.GetChannels error: %w", fmt.Errorf("get channels: some error")),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		messageRepo := &mocks.MessageRepo{}
		replyRepo := &mocks.ReplyRepo{}

		logger := logger.Get(&config.Config{LogLevel: "info"})
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
		messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
		channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
		tt.mock(channelRepo)

		got, err := channelService.GetChannels()
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

func Test_GetChannelsByPage(t *testing.T) {
	data := []model.Channel{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10},
		{ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15}, {ID: 16}, {ID: 17}, {ID: 18}, {ID: 19}, {ID: 20},
	}

	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		input         int
		want          []model.Channel
		expectedError error
	}{
		{
			name: "GetChannelsByPage successful with first page",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(data[:10], nil)
			},
			input: 1,
			want:  data[:10],
		},
		{
			name: "GetChannelsByPage successful with second page",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 10).Return(data[9:], nil)
			},
			input: 2,
			want:  data[9:],
		},
		{
			name: "GetChannelsByPage failed with not found channels by page",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 4030).Return(nil, nil)
			},
			input:         404,
			expectedError: service.ErrChannelsNotFound,
		},
		{name: "GetChannelsByPage failed with some store error",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(nil, fmt.Errorf("get channels by page: some error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Channel] Service.GetChannelsByPage error: %w",
				fmt.Errorf("get channels by page: some error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		messageRepo := &mocks.MessageRepo{}
		replyRepo := &mocks.ReplyRepo{}

		logger := logger.Get(&config.Config{LogLevel: "info"})
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
		messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
		channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
		tt.mock(channelRepo)

		got, err := channelService.GetChannelsByPage(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

func Test_GetChannelByName(t *testing.T) {
	channel := &model.Channel{ID: 1, Name: "test", Title: "test", ImageURL: "test.jpg"}

	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		input         string
		want          *model.Channel
		expectedError error
	}{
		{
			name: "GetChannelByName successful",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(channel, nil)
			},
			input: "test",
			want:  channel,
		},
		{
			name: "GetChannelByName failed with not found channel",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, nil)
			},
			input:         "test",
			expectedError: service.ErrChannelNotFound,
		},
		{
			name: "GetChannelByName failed with some store error",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, fmt.Errorf("get channel by name: some error"))
			},
			input: "test",
			expectedError: fmt.Errorf(
				"[Channel] Service.GetChannelByName error: %w",
				fmt.Errorf("get channel by name: some error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		messageRepo := &mocks.MessageRepo{}
		replyRepo := &mocks.ReplyRepo{}

		logger := logger.Get(&config.Config{LogLevel: "info"})
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
		messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
		channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
		tt.mock(channelRepo)

		got, err := channelService.GetChannelByName(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

func Test_GetChannelStats(t *testing.T) {
	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		input         int
		want          *model.Stat
		wantErr       bool
		err           error
		expectedError error
	}{
		{
			name: "GetChannelStats successful",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelStats", 1).Return(&model.Stat{MessagesCount: 1, RepliesCount: 12}, nil)
			},
			input: 1,
			want:  &model.Stat{MessagesCount: 1, RepliesCount: 12},
		},
		{
			name: "GetChannelStats failed with some store error",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelStats", 1).Return(nil, fmt.Errorf("get channel statistic: some error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Channel] Service.GetChannelStats error: %w",
				fmt.Errorf("get channel statistic: some error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		messageRepo := &mocks.MessageRepo{}
		replyRepo := &mocks.ReplyRepo{}

		logger := logger.Get(&config.Config{LogLevel: "info"})
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
		messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
		channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
		tt.mock(channelRepo)

		got, err := channelService.GetChannelStats(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}
