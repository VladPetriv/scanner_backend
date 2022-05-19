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

func TestReplieService_GetReplies(t *testing.T) {
	replies := []model.Replie{
		{ID: 1, MessageID: 1, UserID: 1, Title: "test1"},
		{ID: 2, MessageID: 1, UserID: 2, Title: "test2"},
	}

	tests := []struct {
		name    string
		mock    func(replieRepo *mocks.ReplieRepo)
		want    []model.Replie
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Replies found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplies").Return(replies, nil)
			},
			want: replies,
		},
		{
			name: "Error: [Replies not found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplies").Return(nil, nil)
			},
			wantErr: true,
			err:     errors.New("replies not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplies").Return(nil, errors.New("error while getting replies: some error"))
			},
			wantErr: true,
			err:     errors.New("[Replie] Service.GetReplies error: error while getting replies: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replieRepo := &mocks.ReplieRepo{}
		replieService := service.NewReplieDBService(&store.Store{Replie: replieRepo})
		tt.mock(replieRepo)

		got, err := replieService.GetReplies()
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

func TestReplieService_GetReplie(t *testing.T) {
	replie := &model.Replie{ID: 1, MessageID: 1, UserID: 1, Title: "test"}

	tests := []struct {
		name    string
		mock    func(replieRepo *mocks.ReplieRepo)
		input   int
		want    *model.Replie
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Replie found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplie", 1).Return(replie, nil)
			},
			input: 1,
			want:  replie,
		},
		{
			name: "Error: [Replie not found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplie", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("replie not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplie", 1).Return(nil, errors.New("error while getting replie: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Replie] Service.GetReplie error: error while getting replie: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replieRepo := &mocks.ReplieRepo{}
		replieService := service.NewReplieDBService(&store.Store{Replie: replieRepo})
		tt.mock(replieRepo)

		got, err := replieService.GetReplie(tt.input)
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

func TestReplieService_GetReplieByName(t *testing.T) {
	replie := &model.Replie{ID: 1, UserID: 1, MessageID: 1, Title: "test"}

	tests := []struct {
		name    string
		mock    func(replieRepo *mocks.ReplieRepo)
		input   string
		want    *model.Replie
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Replie found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", "test").Return(replie, nil)
			},
			input: "test",
			want:  replie,
		},
		{
			name: "Error: [Replie not found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", "test").Return(nil, nil)
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("replie not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", "test").Return(nil, errors.New("error while getting replie: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[Replie] Service.GetReplieByName error: error while getting replie: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replieRepo := &mocks.ReplieRepo{}
		replieService := service.NewReplieDBService(&store.Store{Replie: replieRepo})
		tt.mock(replieRepo)

		got, err := replieService.GetReplieByName(tt.input)
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

func TestReplieService_GetFullRepliesByMessageID(t *testing.T) {
	replies := []model.FullReplie{
		{ID: 1, UserID: 1, Title: "test1", FullName: "test test1", PhotoURL: "test1.jpg"},
		{ID: 2, UserID: 2, Title: "test2", FullName: "test test2", PhotoURL: "test2.jpg"},
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
				replieRepo.On("GetFullRepliesByMessageID", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("replies not found"),
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
