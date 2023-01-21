package service_test

import (
	"fmt"
	"testing"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_CreateWebUser(t *testing.T) {
	input := &model.WebUser{Email: "test@test.com", Password: "test"}

	tests := []struct {
		name          string
		mock          func(webUserRepo *mocks.WebUserRepo)
		input         *model.WebUser
		expectedError error
	}{
		{
			name: "CreateWebUser successful",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", input).Return(nil)
			},
			input: input,
		},
		{
			name: "CreteWebUser failed with some store error",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", input).Return(fmt.Errorf("create web user: some store error"))
			},
			input: input,
			expectedError: fmt.Errorf(
				"[WebUser] Service.CreateWebUser error: %w",
				fmt.Errorf("create web user: some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running %s", tt.name)

		webUserRepo := &mocks.WebUserRepo{}
		webUserService := service.NewWebUserService(&store.Store{WebUser: webUserRepo})
		tt.mock(webUserRepo)

		err := webUserService.CreateWebUser(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
		}

		webUserRepo.AssertExpectations(t)
	}
}

func Test_GetWebUserByID(t *testing.T) {
	user := &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"}

	tests := []struct {
		name          string
		mock          func(webUserRepo *mocks.WebUserRepo)
		input         int
		want          *model.WebUser
		expectedError error
	}{
		{
			name: "GetWebUserByID successful",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByID", 1).Return(user, nil)
			},
			input: 1,
			want:  user,
		},
		{
			name: "GetWebUserByID failed with not found user",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByID", 1).Return(nil, nil)
			},
			input:         1,
			expectedError: service.ErrWebUserNotFound,
		},
		{
			name: "GetWebUserByID failed with some store error",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByID", 1).Return(nil, fmt.Errorf("get web user by id: some store error"))
			},
			input: 1,
			expectedError: fmt.Errorf(
				"[WebUser] Service.GetWebUserByID error: %w",
				fmt.Errorf("get web user by id: some store error"),
			),
		},
	}
	for _, tt := range tests {
		t.Logf("running %s", tt.name)

		webUserRepo := &mocks.WebUserRepo{}
		webUserService := service.NewWebUserService(&store.Store{WebUser: webUserRepo})
		tt.mock(webUserRepo)

		got, err := webUserService.GetWebUserByID(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		webUserRepo.AssertExpectations(t)
	}
}

func Test_GetWebUserByEmail(t *testing.T) {
	user := &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"}

	tests := []struct {
		name          string
		mock          func(webUserRepo *mocks.WebUserRepo)
		input         string
		want          *model.WebUser
		expectedError error
	}{
		{
			name: "GetWebUserByEmail successful",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(user, nil)
			},
			input: "test@test.com",
			want:  user,
		},
		{
			name: "GetWebUserByEmail failed with not found user",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, nil)
			},
			input:         "test@test.com",
			expectedError: service.ErrWebUserNotFound,
		},
		{
			name: "GetWebUserByEmail failed with some store error",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").
					Return(nil, fmt.Errorf("get web user by email: some store error"))
			},
			input: "test@test.com",
			expectedError: fmt.Errorf(
				"[WebUser] Service.GetWebUserByEmail error: %w",
				fmt.Errorf("get web user by email: some store error"),
			),
		},
	}

	for _, tt := range tests {
		t.Logf("running %s", tt.name)

		webUserRepo := &mocks.WebUserRepo{}
		webUserService := service.NewWebUserService(&store.Store{WebUser: webUserRepo})
		tt.mock(webUserRepo)

		got, err := webUserService.GetWebUserByEmail(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		webUserRepo.AssertExpectations(t)
	}
}
