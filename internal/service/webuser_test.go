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

func TestWebUserService_GetWebUser(t *testing.T) {
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
				webUserRepo.On("GetWebUser", 1).Return(user, nil)
			},
			input: 1,
			want:  user,
		},
		{
			name: "Error: [User not found]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUser", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("web user not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUser", 1).Return(nil, errors.New("error while getting web user: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[WebUser] Service.GetWebUser error: error while getting web user: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running %s", tt.name)

		webUserRepo := &mocks.WebUserRepo{}
		webUserService := service.NewWebUserDbService(&store.Store{WebUser: webUserRepo})
		tt.mock(webUserRepo)

		got, err := webUserService.GetWebUser(tt.input)
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

func TestWebUserService_GetWebUserByEmail(t *testing.T) {
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
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, nil)
			},
			input:   "test@test.com",
			wantErr: true,
			err:     errors.New("web user not found"),
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

func TestWebUserService_CreateWebUser(t *testing.T) {
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
