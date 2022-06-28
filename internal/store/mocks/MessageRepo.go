// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/VladPetriv/scanner_backend/internal/model"

// MessageRepo is an autogenerated mock type for the MessageRepo type
type MessageRepo struct {
	mock.Mock
}

// GetFullMessageByMessageID provides a mock function with given fields: ID
func (_m *MessageRepo) GetFullMessageByMessageID(ID int) (*model.FullMessage, error) {
	ret := _m.Called(ID)

	var r0 *model.FullMessage
	if rf, ok := ret.Get(0).(func(int) *model.FullMessage); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FullMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFullMessages provides a mock function with given fields: page
func (_m *MessageRepo) GetFullMessages(page int) ([]model.FullMessage, error) {
	ret := _m.Called(page)

	var r0 []model.FullMessage
	if rf, ok := ret.Get(0).(func(int) []model.FullMessage); ok {
		r0 = rf(page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FullMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFullMessagesByChannelID provides a mock function with given fields: ID, limit, page
func (_m *MessageRepo) GetFullMessagesByChannelID(ID int, limit int, page int) ([]model.FullMessage, error) {
	ret := _m.Called(ID, limit, page)

	var r0 []model.FullMessage
	if rf, ok := ret.Get(0).(func(int, int, int) []model.FullMessage); ok {
		r0 = rf(ID, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FullMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, int) error); ok {
		r1 = rf(ID, limit, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFullMessagesByUserID provides a mock function with given fields: ID
func (_m *MessageRepo) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
	ret := _m.Called(ID)

	var r0 []model.FullMessage
	if rf, ok := ret.Get(0).(func(int) []model.FullMessage); ok {
		r0 = rf(ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FullMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessagesLength provides a mock function with given fields:
func (_m *MessageRepo) GetMessagesLength() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessagesLengthByChannelID provides a mock function with given fields: ID
func (_m *MessageRepo) GetMessagesLengthByChannelID(ID int) (int, error) {
	ret := _m.Called(ID)

	var r0 int
	if rf, ok := ret.Get(0).(func(int) int); ok {
		r0 = rf(ID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
