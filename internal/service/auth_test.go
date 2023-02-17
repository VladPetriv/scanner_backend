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

func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	input := &model.WebUser{Email: "test@test.com", Password: "test"}

	tests := []struct {
		name          string
		mock          func(webUserRepo *mocks.WebUserRepo)
		input         *model.WebUser
		expectedError error
	}{
		{
			name: "Register successful",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, nil)
				webUserRepo.On("CreateWebUser", input).Return(nil)
			},
			input: input,
		},
		{
			name: "Register failed with existed user",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(input, nil)
			},
			input:         input,
			expectedError: service.ErrWebUserIsExist,
		},
		{
			name: "Register failed with some store error when get user by email",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("some store error"))
			},
			input: &model.WebUser{Email: "test@test.com", Password: "test"},
			expectedError: fmt.Errorf(
				"[Register]: %w",
				fmt.Errorf("get web user by email from db: %w", fmt.Errorf("some store error")),
			),
		},
		{
			name: "Register failed with some store error when create user",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, nil)
				webUserRepo.On("CreateWebUser", input).Return(fmt.Errorf("some store error"))
			},
			input: input,
			expectedError: fmt.Errorf(
				"[Register]: %w",
				fmt.Errorf("create web user in db: %w", fmt.Errorf("some store error")),
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
			authService := service.NewAuthService(webUserService, logger)
			tt.mock(webUserRepo)

			err := authService.Register(tt.input)
			assert.Equal(t, tt.expectedError, err)

			webUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	returned := &model.WebUser{
		Email:    "test@test.com",
		Password: "$2a$14$EutZAenAn0GJ223ZMKX/h.WSz8pMuejC0D1QerS5160ibJqjG1Eve",
	}

	tests := []struct {
		name          string
		mock          func(webUserRepo *mocks.WebUserRepo)
		input         *model.WebUser
		want          string
		expectedError error
	}{
		{
			name: "Login successful",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(returned, nil)
			},
			want: "test@test.com",
			input: &model.WebUser{
				Email:    "test@test.com",
				Password: "test",
			},
		},
		{
			name: "Login failed with not found user",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, nil)
			},
			input: &model.WebUser{
				Email: "test@test.com",
			},
			expectedError: service.ErrWebUserNotFound,
		},
		{
			name: "Login failed with incorrect password",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(returned, nil)
			},
			input: &model.WebUser{
				Email:    "test@test.com",
				Password: "test1",
			},
			expectedError: service.ErrIncorrectPassword,
		},
		{
			name: "Login failed with some store error when get user",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("some store error"))
			},
			input: &model.WebUser{
				Email: "test@test.com",
			},
			expectedError: fmt.Errorf(
				"[Login]: %w",
				fmt.Errorf("get web user by email from db: %w", fmt.Errorf("some store error")),
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
			authService := service.NewAuthService(webUserService, logger)
			tt.mock(webUserRepo)

			got, err := authService.Login(tt.input.Email, tt.input.Password)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.want, got)

			webUserRepo.AssertExpectations(t)
		})
	}
}
