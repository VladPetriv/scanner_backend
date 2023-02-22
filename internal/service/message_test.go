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

func TestMessageService_CreateMessage(t *testing.T) {
	t.Parallel()

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
				messageRepo.On("CreateMessage", &model.DBMessage{
					ChannelID:  1,
					UserID:     1,
					Title:      "test",
					MessageURL: "test.url",
					ImageURL:   "test.jpg",
				}).Return(1, nil)
			},
			input: &model.DBMessage{
				ChannelID:  1,
				UserID:     1,
				Title:      "test",
				MessageURL: "test.url",
				ImageURL:   "test.jpg",
			},
			want: 1,
		},
		{
			name: "CreateMessage failed with existed message",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(&model.DBMessage{
					ChannelID:  1,
					UserID:     1,
					Title:      "test",
					MessageURL: "test.url",
					ImageURL:   "test.jpg",
				}, nil)
			},
			input: &model.DBMessage{
				ChannelID:  1,
				UserID:     1,
				Title:      "test",
				MessageURL: "test.url",
				ImageURL:   "test.jpg",
			},
			expectedError: service.ErrMessageExists,
		},
		{
			name: "CreateMessage failed with some store error when get message by title",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(nil, fmt.Errorf("some store error"))
			},
			input: &model.DBMessage{
				ChannelID:  1,
				UserID:     1,
				Title:      "test",
				MessageURL: "test.url",
				ImageURL:   "test.jpg",
			},
			expectedError: fmt.Errorf(
				"get message by title from db: %w",
				fmt.Errorf("some store error"),
			),
		},
		{
			name: "CreateMessage failed with some store error",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByTitle", "test").Return(nil, nil)
				messageRepo.On("CreateMessage", &model.DBMessage{
					ChannelID:  1,
					UserID:     1,
					Title:      "test",
					MessageURL: "test.url",
					ImageURL:   "test.jpg",
				},
				).Return(0, fmt.Errorf("some store error"))
			},
			input: &model.DBMessage{
				ChannelID:  1,
				UserID:     1,
				Title:      "test",
				MessageURL: "test.url",
				ImageURL:   "test.jpg",
			},
			expectedError: fmt.Errorf(
				"create message in db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo)

			got, err := messageService.CreateMessage(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_GetMessagesCount(t *testing.T) {
	t.Parallel()

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
				messageRepo.On("GetMessagesCount").Return(0, fmt.Errorf("some store error"))
			},
			expectedError: fmt.Errorf(
				"get messages count from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo)

			got, err := messageService.GetMessagesCount()
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_GetMessagesCountByChannelID(t *testing.T) {
	t.Parallel()

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
					Return(0, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get messages count by channel id from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo)

			got, err := messageService.GetMessagesCountByChannelID(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_GetFullMessagesByChannelIDAndPage(t *testing.T) {
	t.Parallel()

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
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get full messages by channel id and page from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo)

			got, err := messageService.GetFullMessagesByChannelIDAndPage(tt.input, 1)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_GetFullMessagesByUserID(t *testing.T) {
	t.Parallel()

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
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get full messages by user id from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo)

			got, err := messageService.GetFullMessagesByUserID(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessagesByMessageID(t *testing.T) {
	t.Parallel()

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
				messageRepo.On("GetFullMessageByID", 1).Return(&model.FullMessage{
					ID:        1,
					Title:     "test1",
					ChannelID: 1,
				}, nil)
			},
			input: 1,
			want: &model.FullMessage{
				ID:        1,
				Title:     "test1",
				ChannelID: 1,
			},
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
				messageRepo.On("GetFullMessageByID", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get full message by message id from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo)

			got, err := messageService.GetFullMessageByMessageID(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_ProcessMessagePager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo, replyRepo *mocks.ReplyRepo)
		input         int
		expectedError error
		want          *service.LoadMessageOutput
	}{
		{
			name: "ProcessMessagePage successful",
			mock: func(messageRepo *mocks.MessageRepo, replyRepo *mocks.ReplyRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(&model.FullMessage{
					ID:    1,
					Title: "test",
				}, nil)
				replyRepo.On("GetFullRepliesByMessageID", 1).Return([]model.FullReply{
					{
						ID:    1,
						Title: "test 1",
					},
					{
						ID:    2,
						Title: "test 2",
					},
				}, nil)
			},
			input: 1,
			want: &service.LoadMessageOutput{
				Message: &model.FullMessage{
					ID:    1,
					Title: "test",
					Replies: []model.FullReply{
						{
							ID:    1,
							Title: "test 1",
						},
						{
							ID:    2,
							Title: "test 2",
						},
					},
					RepliesCount: 2,
				},
			},
		},
		{
			name: "ProcessMessagePage failed with not found message",
			mock: func(messageRepo *mocks.MessageRepo, replyRepo *mocks.ReplyRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(nil, nil)
			},
			input: 1,
			want:  &service.LoadMessageOutput{},
		},
		{
			name: "ProcessMessagePage failed with not found replies",
			mock: func(messageRepo *mocks.MessageRepo, replyRepo *mocks.ReplyRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(nil, nil)
			},
			input: 1,
			want:  &service.LoadMessageOutput{},
		},
		{
			name: "ProcessMessagePage failed with some store error when get message by id",
			mock: func(messageRepo *mocks.MessageRepo, replyRepo *mocks.ReplyRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get full message by id from db: %w",
				fmt.Errorf("some store error"),
			),
		},
		{
			name: "ProcessMessagePage failed with some store error when get full replies by message id",
			mock: func(messageRepo *mocks.MessageRepo, replyRepo *mocks.ReplyRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(&model.FullMessage{
					ID:    1,
					Title: "test",
				}, nil)
				replyRepo.On("GetFullRepliesByMessageID", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[ProcessMessagePage]: %w",
				fmt.Errorf(
					"get replies by message id from db: %w",
					fmt.Errorf("some store error"),
				),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo, replyRepo)

			got, err := messageService.ProcessMessagePage(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}

func TestMessageService_ProcessHomePage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(messageRepo *mocks.MessageRepo)
		input         int
		want          *service.LoadHomeOutput
		expectedError error
	}{
		{
			name: "ProcessHomePage successful",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(2, nil)
				messageRepo.On("GetFullMessagesByPage", 0).Return([]model.FullMessage{
					{ID: 1},
					{ID: 2},
				}, nil)
			},
			input: 1,
			want: &service.LoadHomeOutput{
				Messages: []model.FullMessage{
					{ID: 1},
					{ID: 2},
				},
				MessagesCount: 2,
			},
		},
		{
			name: "ProcessHomePage failed with not found messages",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, nil)
			},
			input: 1,
			want:  &service.LoadHomeOutput{},
		},
		{
			name: "ProcessHomePage failed with some store error when get messages count",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get messages count from db: %w",
				fmt.Errorf("some store error"),
			),
		},
		{
			name: "ProcessHomePage failed with some store error when get messages by page",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(1, nil)
				messageRepo.On("GetFullMessagesByPage", 0).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get full messages by page from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messageRepo := &mocks.MessageRepo{}
			replyRepo := &mocks.ReplyRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
			messageService := service.NewMessageService(&store.Store{Message: messageRepo}, logger, replyService)
			tt.mock(messageRepo)

			got, err := messageService.ProcessHomePage(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			messageRepo.AssertExpectations(t)
			replyRepo.AssertExpectations(t)
		})
	}
}
