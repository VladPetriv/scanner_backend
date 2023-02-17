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

func TestReplyService_CreateReply(t *testing.T) {
	t.Parallel()

	replyInput := &model.DBReply{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"}

	tests := []struct {
		name          string
		mock          func(replyRepo *mocks.ReplyRepo)
		input         *model.DBReply
		expectedError error
	}{
		{
			name: "CreateReply successful",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("CreateReply", replyInput).Return(nil)
			},
			input: replyInput,
		},
		{
			name: "CreteReply failed with some store error",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("CreateReply", replyInput).Return(fmt.Errorf("some store error"))
			},
			input: replyInput,
			expectedError: fmt.Errorf(
				"create reply in db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replyRepo := &mocks.ReplyRepo{}

		logger := logger.Get(&config.Config{LogLevel: "info"})
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
		tt.mock(replyRepo)

		err := replyService.CreateReply(tt.input)
		assert.Equal(t, tt.expectedError, err)

		replyRepo.AssertExpectations(t)
	}
}

func Test_GetFullRepliesByMessageID(t *testing.T) {
	replies := []model.FullReply{
		{ID: 1, UserID: 1, Title: "test1", FullName: "test test1", UserImageURL: "test1.jpg"},
		{ID: 2, UserID: 2, Title: "test2", FullName: "test test2", UserImageURL: "test2.jpg"},
	}

	tests := []struct {
		name          string
		mock          func(replyRepo *mocks.ReplyRepo)
		input         int
		want          []model.FullReply
		expectedError error
	}{
		{
			name: "GetFullRepliesByMessageID successful",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("GetFullRepliesByMessageID", 1).Return(replies, nil)
			},
			input: 1,
			want:  replies,
		},
		{
			name: "GetFullRepliesByMessageID failed with not found replies",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("GetFullRepliesByMessageID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrRepliesNotFound,
		},
		{
			name: "GetFullRepliesByMessageID failed with store error",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("GetFullRepliesByMessageID", 1).Return(nil, fmt.Errorf("some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"get replies by message id from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replyRepo := &mocks.ReplyRepo{}

		logger := logger.Get(&config.Config{LogLevel: "info"})
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo}, logger)
		tt.mock(replyRepo)

		got, err := replyService.GetFullRepliesByMessageID(tt.input)
		assert.Equal(t, tt.expectedError, err)
		assert.Equal(t, tt.want, got)

		replyRepo.AssertExpectations(t)
	}
}
