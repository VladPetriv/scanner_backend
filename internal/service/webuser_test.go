package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/service"
	"github.com/VladPetriv/scanner_backend/internal/store"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/stretchr/testify/assert"
)

func Test_GetWebUserByID(t *testing.T) {
	user := &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"}

	tests := []struct {
		name    string
		mock    func(webUserRepo *mocks.WebUserRepo)
		input   int
		want    *model.WebUser
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User found]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByID", 1).Return(user, nil)
			},
			input: 1,
			want:  user,
		},
		{
			name: "Error: [User not found]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByID", 1).Return(nil, pg.ErrWebUserNotFound)
			},
			input:   1,
			wantErr: true,
			err:     fmt.Errorf("[WebUser] Service.GetWebUserByID error: %w", pg.ErrWebUserNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByID", 1).Return(nil, errors.New("error while getting web user: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[WebUser] Service.GetWebUserByID error: error while getting web user: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running %s", tt.name)

		webUserRepo := &mocks.WebUserRepo{}
		webUserService := service.NewWebUserDbService(&store.Store{WebUser: webUserRepo})
		tt.mock(webUserRepo)

		got, err := webUserService.GetWebUserByID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
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
		name    string
		mock    func(webUserRepo *mocks.WebUserRepo)
		input   string
		want    *model.WebUser
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User found]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(user, nil)
			},
			input: "test@test.com",
			want:  user,
		},
		{
			name: "Error: [User not found]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, pg.ErrWebUserNotFound)
			},
			input:   "test@test.com",
			wantErr: true,
			err:     fmt.Errorf("[WebUser] Service.GetWebUserByEmail error: %w", pg.ErrWebUserNotFound),
		},
		{
			name: "Error: [Store error]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, errors.New("error while getting web user by email: some error"))
			},
			input:   "test@test.com",
			wantErr: true,
			err:     errors.New("[WebUser] Service.GetWebUserByEmail error: error while getting web user by email: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running %s", tt.name)

		webUserRepo := &mocks.WebUserRepo{}
		webUserService := service.NewWebUserDbService(&store.Store{WebUser: webUserRepo})
		tt.mock(webUserRepo)

		got, err := webUserService.GetWebUserByEmail(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		webUserRepo.AssertExpectations(t)
	}
}

func Test_CreateWebUser(t *testing.T) {
	input := &model.WebUser{Email: "test@test.com", Password: "test"}

	tests := []struct {
		name    string
		mock    func(webUserRepo *mocks.WebUserRepo)
		input   *model.WebUser
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User created]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", input).Return(1, nil)
			},
			input: input,
		},
		{
			name: "Error: [Store error]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", input).Return(0, errors.New("error while creating web user: some error"))
			},
			input:   input,
			wantErr: true,
			err:     errors.New("[WebUser] Service.CreateWebUser error: error while creating web user: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running %s", tt.name)

		webUserRepo := &mocks.WebUserRepo{}
		webUserService := service.NewWebUserDbService(&store.Store{WebUser: webUserRepo})
		tt.mock(webUserRepo)

		err := webUserService.CreateWebUser(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, nil, err)
		}

		webUserRepo.AssertExpectations(t)
	}
}
