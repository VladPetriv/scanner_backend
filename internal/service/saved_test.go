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

func TestSavedService_GetSavedMessages(t *testing.T) {
	saved := []model.Saved{
		{ID: 1, WebUserID: 1, MessageID: 1},
		{ID: 2, WebUserID: 1, MessageID: 5},
	}

	tests := []struct {
		name    string
		mock    func(savedRepo *mocks.SavedRepo)
		input   int
		want    []model.Saved
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Saved messages found]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(saved, nil)
			},
			input: 1,
			want:  saved,
		},
		{
			name: "Error: [Saved messages not found]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("saved messages not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, errors.New("error while getting saved messages: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Saved] Service.GetSavedMessages error: error while getting saved messages: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		savedRepo := &mocks.SavedRepo{}
		savedService := service.NewSavedDbService(&store.Store{Saved: savedRepo})
		tt.mock(savedRepo)

		got, err := savedService.GetSavedMessages(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		savedRepo.AssertExpectations(t)
	}
}

func TestSavedService_GetSavedMessageByMessageID(t *testing.T) {
	saved := &model.Saved{ID: 1, WebUserID: 1, MessageID: 1}

	tests := []struct {
		name    string
		mock    func(savedRepo *mocks.SavedRepo)
		input   int
		want    *model.Saved
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Saved message found]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByMessageID", 1).Return(saved, nil)
			},
			input: 1,
			want:  saved,
		},
		{
			name: "Error: [Saved message not found]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByMessageID", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("saved message not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByMessageID", 1).Return(nil, errors.New("error while getting saved message: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Saved] Service.GetSavedMessageByMessageID error: error while getting saved message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		savedRepo := &mocks.SavedRepo{}
		savedService := service.NewSavedDbService(&store.Store{Saved: savedRepo})
		tt.mock(savedRepo)

		got, err := savedService.GetSavedMessageByMessageID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		savedRepo.AssertExpectations(t)
	}
}

func TestSavedService_CreateSavedMessage(t *testing.T) {
	input := &model.Saved{WebUserID: 1, MessageID: 1}

	tests := []struct {
		name    string
		mock    func(savedRepo *mocks.SavedRepo)
		input   *model.Saved
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Saved message created]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", input).Return(1, nil)
			},
			input: input,
		},
		{
			name: "Error: [Store error]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", input).Return(0, errors.New("error while creating saved message: some error"))
			},
			input:   input,
			wantErr: true,
			err:     errors.New("[Saved] Service.CreateSavedMessage error: error while creating saved message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		savedRepo := &mocks.SavedRepo{}
		savedService := service.NewSavedDbService(&store.Store{Saved: savedRepo})
		tt.mock(savedRepo)

		err := savedService.CreateSavedMessage(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}

		savedRepo.AssertExpectations(t)
	}
}
