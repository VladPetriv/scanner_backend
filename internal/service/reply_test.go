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

func Test_CreateReply(t *testing.T) {
	t.Parallel()

	replyInput := &model.DBReply{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"}

	tests := []struct {
		name    string
		mock    func(replyRepo *mocks.ReplyRepo)
		input   *model.DBReply
		wantErr bool
		err     error
	}{
		{
			name: "CreateReply successful",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("CreateReply", replyInput).Return(nil)
			},
			input: replyInput,
		},
		{
			name: "CreteReply failed with some sql error",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("CreateReply", replyInput).Return(fmt.Errorf("failed to create reply: some error"))
			},
			input:   replyInput,
			wantErr: true,
			err:     fmt.Errorf("[Reply] Service.CreateReply error: %w", fmt.Errorf("failed to create reply: some error")),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replyRepo := &mocks.ReplyRepo{}
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo})
		tt.mock(replyRepo)

		err := replyService.CreateReply(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.EqualValues(t, tt.err, err)
		} else {
			assert.NoError(t, err)
		}

		replyRepo.AssertExpectations(t)
	}
}

func Test_GetFullRepliesByMessageID(t *testing.T) {
	replies := []model.FullReply{
		{ID: 1, UserID: 1, Title: "test1", FullName: "test test1", UserImageURL: "test1.jpg"},
		{ID: 2, UserID: 2, Title: "test2", FullName: "test test2", UserImageURL: "test2.jpg"},
	}

	tests := []struct {
		name    string
		mock    func(replyRepo *mocks.ReplyRepo)
		input   int
		want    []model.FullReply
		wantErr bool
		err     error
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
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("replies not found"),
		},
		{
			name: "GetFullRepliesByMessageID failed with store error",
			mock: func(replyRepo *mocks.ReplyRepo) {
				replyRepo.On("GetFullRepliesByMessageID", 1).Return(nil, fmt.Errorf("get full replies by message id: some error"))
			},
			input:   1,
			wantErr: true,
			err: fmt.Errorf(
				"[Reply] Service.GetFullRepliesByMessageID error: %w",
				fmt.Errorf("get full replies by message id: some error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replyRepo := &mocks.ReplyRepo{}
		replyService := service.NewReplyService(&store.Store{Reply: replyRepo})
		tt.mock(replyRepo)

		got, err := replyService.GetFullRepliesByMessageID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		replyRepo.AssertExpectations(t)
	}
}
