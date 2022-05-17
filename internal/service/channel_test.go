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

func TestChannelService_GetChannels(t *testing.T) {
	data := []model.Channel{
		{ID: 1, Name: "test1", Title: "test1", PhotoURL: "test1.jpg"},
		{ID: 2, Name: "test2", Title: "test2", PhotoURL: "test2.jpg"},
	}

	tests := []struct {
		name    string
		mock    func(channelRepo *mocks.ChannelRepo)
		want    []model.Channel
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Channels found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannels").Return(data, nil)
			},
			want: data,
		},
		{
			name: "Error: [Channels not found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannels").Return(nil, nil)
			},
			wantErr: true,
			err:     errors.New("channels not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannels").Return(nil, errors.New("error while getting channels: some error"))
			},
			wantErr: true,
			err:     errors.New("[Channel] Service.GetChannels error: error while getting channels: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		channelService := service.NewChannelDBService(&store.Store{Channel: channelRepo})
		tt.mock(channelRepo)

		got, err := channelService.GetChannels()
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

func TestChannelService_GetChannelsByPage(t *testing.T) {
	data := []model.Channel{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10},
		{ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15}, {ID: 16}, {ID: 17}, {ID: 18}, {ID: 19}, {ID: 20},
	}

	tests := []struct {
		name    string
		mock    func(channelRepo *mocks.ChannelRepo)
		input   int
		want    []model.Channel
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Channels found, return 1 page]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(data[:10], nil)
			},
			input: 1,
			want:  data[:10],
		},
		{
			name: "Ok: [Channels found, return 2 page]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 10).Return(data[9:], nil)
			},
			input: 2,
			want:  data[9:],
		},
		{
			name: "Error: [Channels not found, return page 404]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 4030).Return(nil, nil)
			},
			input:   404,
			wantErr: true,
			err:     errors.New("channels not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(nil, errors.New("error while getting channels: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Channel] Service.GetChannelsByPage error: error while getting channels: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		channelServie := service.NewChannelDBService(&store.Store{Channel: channelRepo})
		tt.mock(channelRepo)

		got, err := channelServie.GetChannelsByPage(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

func TestChannelService_GetChannel(t *testing.T) {
	channel := &model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"}

	tests := []struct {
		name    string
		mock    func(channelRepo *mocks.ChannelRepo)
		input   int
		want    *model.Channel
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Channel found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannel", 1).Return(channel, nil)
			},
			input: 1,
			want:  channel,
		},
		{
			name: "Error: [Channel not found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannel", 1).Return(nil, nil)
			},
			input:   1,
			wantErr: true,
			err:     errors.New("channel not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannel", 1).Return(nil, errors.New("error while getting channel: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Channel] Service.GetChannel error: error while getting channel: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		channelService := service.NewChannelDBService(&store.Store{Channel: channelRepo})
		tt.mock(channelRepo)

		got, err := channelService.GetChannel(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

func TestChannelService_GetChannelByName(t *testing.T) {
	channel := &model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"}

	tests := []struct {
		name    string
		mock    func(channelRepo *mocks.ChannelRepo)
		input   string
		want    *model.Channel
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Channel found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(channel, nil)
			},
			input: "test",
			want:  channel,
		},
		{
			name: "Error: [Channel not found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, nil)
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("channel not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, errors.New("error while getting channel: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[Channel] Service.GetChannelByName error: error while getting channel: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		channelService := service.NewChannelDBService(&store.Store{Channel: channelRepo})
		tt.mock(channelRepo)

		got, err := channelService.GetChannelByName(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

func TestChannelServie_GetChannelStats(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(channelRepo *mocks.ChannelRepo)
		input   int
		want    *model.Stat
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Stat found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelStats", 1).Return(&model.Stat{MessagesCount: 1, RepliesCount: 12}, nil)
			},
			input: 1,
			want:  &model.Stat{MessagesCount: 1, RepliesCount: 12},
		},
		{
			name: "Error: [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelStats", 1).Return(nil, errors.New("error while getting channel stat: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Channel] Service.GetChannelStats error: error while getting channel stat: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		channelService := service.NewChannelDBService(&store.Store{Channel: channelRepo})
		tt.mock(channelRepo)

		got, err := channelService.GetChannelStats(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}
