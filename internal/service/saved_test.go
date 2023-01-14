package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_CreateSavedMessage(t *testing.T) {
	input := &model.Saved{WebUserID: 1, MessageID: 1}

	tests := []struct {
		name          string
		mock          func(savedRepo *mocks.SavedRepo)
		input         *model.Saved
		expectedError error
	}{
		{
			name: "CreateSavedMessage successful",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", input).Return(nil)
			},
			input: input,
		},
		{
			name: "CreateSavedMessage failed with some store error",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", input).Return(fmt.Errorf("create saved message: some store error"))
			},
			input: input,
			expectedError: fmt.Errorf(
				"[Saved] Service.CreateSavedMessage error: %w",
				fmt.Errorf("create saved message: some store error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		savedRepo := &mocks.SavedRepo{}
		savedService := service.NewSavedService(&store.Store{Saved: savedRepo})
		tt.mock(savedRepo)

		err := savedService.CreateSavedMessage(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
		}

		savedRepo.AssertExpectations(t)
	}
}

func Test_GetSavedMessages(t *testing.T) {
	saved := []model.Saved{
		{ID: 1, WebUserID: 1, MessageID: 1},
		{ID: 2, WebUserID: 1, MessageID: 5},
	}

	tests := []struct {
		name          string
		mock          func(savedRepo *mocks.SavedRepo)
		input         int
		want          []model.Saved
		expectedError error
	}{
		{
			name: "GetSavedMessages successful",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(saved, nil)
			},
			input: 1,
			want:  saved,
		},
		{
			name: "GetSavedMessages failed with not found messages",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrSavedMessagesNotFound,
		},
		{
			name: "GetSavedMessages failed with some store error",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, fmt.Errorf("get saved messages: some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Saved] Service.GetSavedMessages error: %w",
				fmt.Errorf("get saved messages: some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		savedRepo := &mocks.SavedRepo{}
		savedService := service.NewSavedService(&store.Store{Saved: savedRepo})
		tt.mock(savedRepo)

		got, err := savedService.GetSavedMessages(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		savedRepo.AssertExpectations(t)
	}
}

func Test_GetSavedMessageByMessageID(t *testing.T) {
	saved := &model.Saved{ID: 1, WebUserID: 1, MessageID: 1}

	tests := []struct {
		name          string
		mock          func(savedRepo *mocks.SavedRepo)
		input         int
		want          *model.Saved
		expectedError error
	}{
		{
			name: "GetSavedMessageByMessageID successful",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByID", 1).Return(saved, nil)
			},
			input: 1,
			want:  saved,
		},
		{
			name: "GetSavedMessageByMessageID failed with not found message",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrSavedMessageNotFound,
		},
		{
			name: "GetSavedMessageByMessageID failed with some store erorr",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessageByID", 1).Return(nil, errors.New("get saved message: some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Saved] Service.GetSavedMessageByMessageID error: %w",
				fmt.Errorf("get saved message: some store error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		savedRepo := &mocks.SavedRepo{}
		savedService := service.NewSavedService(&store.Store{Saved: savedRepo})
		tt.mock(savedRepo)

		got, err := savedService.GetSavedMessageByMessageID(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		savedRepo.AssertExpectations(t)
	}
}

func Test_DeleteSavedMessage(t *testing.T) {
	tests := []struct {
		name          string
		mock          func(savedRepo *mocks.SavedRepo)
		input         int
		expectedError error
	}{
		{
			name: "DeleteSavedMessage successful",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("DeleteSavedMessage", 1).Return(nil)
			},
			input: 1,
		},
		{
			name: "DeleteSavedMessage failed with some store error",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("DeleteSavedMessage", 1).Return(fmt.Errorf("delete saved message: some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[Saved] Service.DeleteSavedMessage error: %w",
				fmt.Errorf("delete saved message: some store error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		savedRepo := &mocks.SavedRepo{}
		savedService := service.NewSavedService(&store.Store{Saved: savedRepo})
		tt.mock(savedRepo)

		err := savedService.DeleteSavedMessage(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
		}

		savedRepo.AssertExpectations(t)
	}
}
