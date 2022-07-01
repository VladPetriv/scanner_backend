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

func TestUserService_GetUsers(t *testing.T) {
	users := []model.User{
		{ID: 1, Username: "test1", FullName: "test1 test", ImageURL: "test1.jpg"},
		{ID: 2, Username: "test2", FullName: "test2 test", ImageURL: "test2.jpg"},
		{ID: 3, Username: "test3", FullName: "test3 test", ImageURL: "test3.jpg"},
	}

	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		want    []model.User
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Users found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUsers").Return(users, nil)
			},
			want: users,
		},
		{
			name: "Error: [Users not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUsers").Return(nil, nil)
			},
			wantErr: true,
			err:     errors.New("users not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUsers").Return(nil, errors.New("error while getting users: some error"))
			},
			wantErr: true,
			err:     errors.New("[User] Service.GetUsers error: error while getting users: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUsers()
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}

func TestUserService_GetUserByUsername(t *testing.T) {
	user := &model.User{
		ID:       1,
		Username: "test1",
		FullName: "test1 test",
		ImageURL: "test1.jpg",
	}

	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		input   string
		want    *model.User
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Users found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(user, nil)
			},
			input: "test",
			want:  user,
		},
		{
			name: "Error: [Users not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, nil)
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("user not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, errors.New("error while getting users: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[User] Service.GetUserByUsername error: error while getting users: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUserByUsername(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	user := &model.User{
		ID:       1,
		Username: "test1",
		FullName: "test1 test",
		ImageURL: "test1.jpg",
	}

	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		input   int
		want    *model.User
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Users found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(user, nil)
			},
			input: 1,
			want:  user,
		},
		{
			name: "Error: [Users not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("user not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, errors.New("error while getting users: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[User] Service.GetUserByID error: error while getting users: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUserByID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}
