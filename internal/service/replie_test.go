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

func Test_CreateReplie(t *testing.T) {
	replieInput := &model.DBReplie{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"}

	tests := []struct {
		name    string
		mock    func(replieRepo *mocks.ReplieRepo)
		input   *model.DBReplie
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Replie created]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("CreateReplie", replieInput).Return(nil)
			},
			input: replieInput,
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("CreateReplie", replieInput).Return(errors.New("failed to create replie: some error"))
			},
			input:   replieInput,
			wantErr: true,
			err:     errors.New("[Replie] Service.CreateReplie error: failed to create replie: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replieRepo := &mocks.ReplieRepo{}
		replieService := service.NewReplieDBService(&store.Store{Replie: replieRepo})
		tt.mock(replieRepo)

		err := replieService.CreateReplie(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.EqualValues(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}

		replieRepo.AssertExpectations(t)
	}
}

func Test_GetFullRepliesByMessageID(t *testing.T) {
	replies := []model.FullReplie{
		{ID: 1, UserID: 1, Title: "test1", FullName: "test test1", UserImageURL: "test1.jpg"},
		{ID: 2, UserID: 2, Title: "test2", FullName: "test test2", UserImageURL: "test2.jpg"},
	}

	tests := []struct {
		name    string
		mock    func(replieRepo *mocks.ReplieRepo)
		input   int
		want    []model.FullReplie
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Replies found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetFullRepliesByMessageID", 1).Return(replies, nil)
			},
			input: 1,
			want:  replies,
		},
		{
			name: "Error: [Replies not found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetFullRepliesByMessageID", 1).Return(nil, pg.ErrRepliesNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[Replie] Service.GetFullRepliesByMessageID error: %w", pg.ErrRepliesNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetFullRepliesByMessageID", 1).Return(nil, errors.New("error while getting full replies by message ID: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Replie] Service.GetFullRepliesByMessageID error: error while getting full replies by message ID: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replieRepo := &mocks.ReplieRepo{}
		replieService := service.NewReplieDBService(&store.Store{Replie: replieRepo})
		tt.mock(replieRepo)

		got, err := replieService.GetFullRepliesByMessageID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		replieRepo.AssertExpectations(t)
	}
}
