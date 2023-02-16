package service_test

import (
	"fmt"
	"testing"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/VladPetriv/scanner_backend/pkg/config"
	"github.com/VladPetriv/scanner_backend/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_CreateWebUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(webUserRepo *mocks.WebUserRepo)
		input         *model.WebUser
		expectedError error
	}{
		{
			name: "CreateWebUser successful",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", &model.WebUser{Email: "test@test.com", Password: "test"}).Return(nil)
			},
			input: &model.WebUser{Email: "test@test.com", Password: "test"},
		},
		{
			name: "CreteWebUser failed with some store error",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", &model.WebUser{
					Email:    "test@test.com",
					Password: "test"},
				).Return(fmt.Errorf("some store error"))
			},
			input: &model.WebUser{Email: "test@test.com", Password: "test"},
			expectedError: fmt.Errorf(
				"create web user in db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			webUserRepo := &mocks.WebUserRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			webUserService := service.NewWebUserService(&store.Store{WebUser: webUserRepo}, logger)
			tt.mock(webUserRepo)

			err := webUserService.CreateWebUser(tt.input)
			assert.Equal(t, tt.expectedError, err)

			webUserRepo.AssertExpectations(t)
		})
	}
}

func Test_GetWebUserByEmail(t *testing.T) {
	t.Parallel()

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
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(&model.WebUser{
					ID:       1,
					Email:    "test@test.com",
					Password: "test",
				}, nil)
			},
			input: "test@test.com",
			want:  &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"},
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
					Return(nil, fmt.Errorf("some store error"))
			},
			input: "test@test.com",
			expectedError: fmt.Errorf(
				"get web user by email from db: %w",
				fmt.Errorf("some store error"),
			),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			webUserRepo := &mocks.WebUserRepo{}

			logger := logger.Get(&config.Config{LogLevel: "info"})
			webUserService := service.NewWebUserService(&store.Store{WebUser: webUserRepo}, logger)
			tt.mock(webUserRepo)

			got, err := webUserService.GetWebUserByEmail(tt.input)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			webUserRepo.AssertExpectations(t)
		})
	}
}
