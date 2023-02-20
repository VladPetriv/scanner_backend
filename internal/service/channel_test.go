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

func TestChannelService_CreateChannel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		input         *model.DBChannel
		expectedError error
	}{
		{
			name: "CreateChannel successful",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, nil)
				channelRepo.On("CreateChannel", &model.DBChannel{
					Name:     "test",
					Title:    "test T",
					ImageURL: "test.jpg",
				}).Return(nil)
			},
			input: &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
		},
		{
			name: "CreateChannel failed with existed channel",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(&model.Channel{Name: "test"}, nil)
			},
			input:         &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
			expectedError: service.ErrChannelExists,
		},
		{
			name: "CreateChannel failed with some store error when get channel by name",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, fmt.Errorf("some store error"))
			},
			input: &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
			expectedError: fmt.Errorf(
				"[CreateChannel]: %w",
				fmt.Errorf(
					"get channel by name from db: %w",
					fmt.Errorf("some store error"),
				),
			),
		},
		{
			name: "CreateChannel failed with some store error when create channel",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, nil)
				channelRepo.On("CreateChannel", &model.DBChannel{Name: "test",
					Title:    "test T",
					ImageURL: "test.jpg",
				}).Return(fmt.Errorf("some store error"))
			},
			input: &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
			expectedError: fmt.Errorf(
				"create channel in db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channelRepo := &mocks.ChannelRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
			tt.mock(channelRepo)

			err := channelService.CreateChannel(tt.input)
			assert.Equal(t, tt.expectedError, err)

			channelRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestChannelService_GetChannels(t *testing.T) {
	t.Parallel()

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
				channelRepo.On("GetChannels").Return(nil, fmt.Errorf("some store error"))
			},
			expectedError: fmt.Errorf(
				"get channels from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channelRepo := &mocks.ChannelRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
			tt.mock(channelRepo)

			got, err := channelService.GetChannels()
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			channelRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestChannelService_GetChannelsByPage(t *testing.T) {
	t.Parallel()

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
				channelRepo.On("GetChannelsByPage", 0).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get channels by page from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channelRepo := &mocks.ChannelRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
			tt.mock(channelRepo)

			got, err := channelService.GetChannelsByPage(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			channelRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestChannelService_GetChannelByName(t *testing.T) {
	t.Parallel()

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
				channelRepo.On("GetChannelByName", "test").Return(&model.Channel{
					ID:       1,
					Name:     "test",
					Title:    "test",
					ImageURL: "test.jpg",
				}, nil)
			},
			input: "test",
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", ImageURL: "test.jpg"},
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
				channelRepo.On("GetChannelByName", "test").Return(nil, fmt.Errorf("some store error"))
			},
			input: "test",
			expectedError: fmt.Errorf(
				"get channel by name from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channelRepo := &mocks.ChannelRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
			tt.mock(channelRepo)

			got, err := channelService.GetChannelByName(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			channelRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestChannelService_GetChannelStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		input         int
		want          *model.Stat
		wantErr       bool
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
				channelRepo.On("GetChannelStats", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get channel statistic from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channelRepo := &mocks.ChannelRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
			tt.mock(channelRepo)

			got, err := channelService.GetChannelStats(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			channelRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestChannelService_ProcessChannelPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		mock             func(channelRepo *mocks.ChannelRepo, messageRepo *mocks.MessageRepo)
		inputChannelName string
		inputPage        int
		expectedError    error
		want             *service.LoadChannelOutput
	}{
		{
			name: "ProcessChannelPage successful",
			mock: func(channelRepo *mocks.ChannelRepo, messageRepo *mocks.MessageRepo) {
				channelRepo.On("GetChannelByName", "test").Return(&model.Channel{
					ID:   1,
					Name: "test",
				}, nil)

				messageRepo.On("GetMessagesCountByChannelID", 1).Return(2, nil)
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return([]model.FullMessage{
					{ID: 1, ChannelID: 1},
					{ID: 2, ChannelID: 1},
				}, nil)
			},
			inputChannelName: "test",
			inputPage:        1,
			want: &service.LoadChannelOutput{
				Channel: model.Channel{
					ID:   1,
					Name: "test",
				},
				Messages: []model.FullMessage{
					{ID: 1, ChannelID: 1},
					{ID: 2, ChannelID: 1},
				},
				MessagesCount: 2,
			},
		},
		{
			name: "ProcessChannelPage failed with not found channel",
			mock: func(channelRepo *mocks.ChannelRepo, messageRepo *mocks.MessageRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, nil)
			},
			inputChannelName: "test",
			inputPage:        1,
			want:             &service.LoadChannelOutput{},
		},
		{
			name: "ProcessChannelPage failed with not found messages count",
			mock: func(channelRepo *mocks.ChannelRepo, messageRepo *mocks.MessageRepo) {
				channelRepo.On("GetChannelByName", "test").Return(&model.Channel{
					ID:   1,
					Name: "test",
				}, nil)
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(0, nil)
			},
			inputChannelName: "test",
			inputPage:        1,
			want: &service.LoadChannelOutput{
				Channel: model.Channel{
					ID:   1,
					Name: "test",
				},
			},
		},
		{
			name: "ProcessChannelPage failed with some store error when get channel by name",
			mock: func(channelRepo *mocks.ChannelRepo, messageRepo *mocks.MessageRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, fmt.Errorf("some store error"))
			},
			inputChannelName: "test",
			inputPage:        1,
			expectedError: fmt.Errorf(
				"[ProcessChannelPage]: %w",
				fmt.Errorf(
					"get channel by name from db: %w",
					fmt.Errorf("some store error"),
				),
			),
		},
		{
			name: "ProcessChannelPage failed with some store error when get messages count",
			mock: func(channelRepo *mocks.ChannelRepo, messageRepo *mocks.MessageRepo) {
				channelRepo.On("GetChannelByName", "test").Return(&model.Channel{
					ID:   1,
					Name: "test",
				}, nil)
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(0, fmt.Errorf("some store error"))
			},
			inputChannelName: "test",
			inputPage:        1,
			expectedError: fmt.Errorf(
				"[ProcessChannelPage]: %w",
				fmt.Errorf(
					"get messages count by channel id from db: %w",
					fmt.Errorf("some store error"),
				),
			),
		},
		{
			name: "ProcessChannelPage failed with some store error when get messages",
			mock: func(channelRepo *mocks.ChannelRepo, messageRepo *mocks.MessageRepo) {
				channelRepo.On("GetChannelByName", "test").Return(&model.Channel{
					ID:   1,
					Name: "test",
				}, nil)
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(1, nil)
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(nil, fmt.Errorf("some store error"))
			},
			inputChannelName: "test",
			inputPage:        1,
			expectedError: fmt.Errorf(
				"[ProcessChannelPage]: %w",
				fmt.Errorf(
					"get full messages by channel id and page from db: %w",
					fmt.Errorf("some store error"),
				),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channelRepo := &mocks.ChannelRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
			tt.mock(channelRepo, messageRepo)

			got, err := channelService.ProcessChannelPage(tt.inputChannelName, tt.inputPage)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			channelRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestChannelService_ProcessChannelsPage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(channelRepo *mocks.ChannelRepo)
		input         int
		expectedError error
		want          *service.LoadChannelsOutput
	}{
		{
			name: "ProcessChannelPage successful",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return([]model.Channel{
					{ID: 1},
					{ID: 2},
				}, nil)
				channelRepo.On("GetChannelStats", 1).Return(&model.Stat{
					MessagesCount: 1,
					RepliesCount:  2,
				}, nil)
				channelRepo.On("GetChannelStats", 2).Return(&model.Stat{
					MessagesCount: 3,
					RepliesCount:  2,
				}, nil)
			},
			input: 1,
			want: &service.LoadChannelsOutput{
				Channels: []model.Channel{
					{ID: 1, Stats: model.Stat{MessagesCount: 1, RepliesCount: 2}},
					{ID: 2, Stats: model.Stat{MessagesCount: 3, RepliesCount: 2}},
				},
			},
		},
		{
			name: "ProcessChannelPage failed with not found channels",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(nil, nil)
			},
			input: 1,
			want:  &service.LoadChannelsOutput{},
		},
		{
			name: "ProcessChannelPage failed some store when get channels by page",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[ProcessChannelsPage]: %w",
				fmt.Errorf(
					"get channels by page from db: %w",
					fmt.Errorf("some store error"),
				),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channelRepo := &mocks.ChannelRepo{}
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			channelService := service.NewChannelService(&store.Store{Channel: channelRepo}, logger, messageService)
			tt.mock(channelRepo)

			got, err := channelService.ProcessChannelsPage(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			channelRepo.AssertExpectations(t)
			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}
