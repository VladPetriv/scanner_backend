// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	model "github.com/VladPetriv/scanner_backend/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// MessageRepo is an autogenerated mock type for the MessageRepo type
type MessageRepo struct {
	mock.Mock
}

// CreateMessage provides a mock function with given fields: message
func (_m *MessageRepo) CreateMessage(message *model.DBMessage) (int, error) {
	ret := _m.Called(message)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.DBMessage) (int, error)); ok {
		return rf(message)
	}
	if rf, ok := ret.Get(0).(func(*model.DBMessage) int); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*model.DBMessage) error); ok {
		r1 = rf(message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFullMessageByID provides a mock function with given fields: id
func (_m *MessageRepo) GetFullMessageByID(id int) (*model.FullMessage, error) {
	ret := _m.Called(id)

	var r0 *model.FullMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*model.FullMessage, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *model.FullMessage); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FullMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFullMessagesByChannelIDAndPage provides a mock function with given fields: id, page
func (_m *MessageRepo) GetFullMessagesByChannelIDAndPage(id int, page int) ([]model.FullMessage, error) {
	ret := _m.Called(id, page)

	var r0 []model.FullMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]model.FullMessage, error)); ok {
		return rf(id, page)
	}
	if rf, ok := ret.Get(0).(func(int, int) []model.FullMessage); ok {
		r0 = rf(id, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FullMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(id, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFullMessagesByPage provides a mock function with given fields: page
func (_m *MessageRepo) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	ret := _m.Called(page)

	var r0 []model.FullMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]model.FullMessage, error)); ok {
		return rf(page)
	}
	if rf, ok := ret.Get(0).(func(int) []model.FullMessage); ok {
		r0 = rf(page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FullMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFullMessagesByUserID provides a mock function with given fields: id
func (_m *MessageRepo) GetFullMessagesByUserID(id int) ([]model.FullMessage, error) {
	ret := _m.Called(id)

	var r0 []model.FullMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]model.FullMessage, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) []model.FullMessage); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FullMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessageByTitle provides a mock function with given fields: title
func (_m *MessageRepo) GetMessageByTitle(title string) (*model.DBMessage, error) {
	ret := _m.Called(title)

	var r0 *model.DBMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.DBMessage, error)); ok {
		return rf(title)
	}
	if rf, ok := ret.Get(0).(func(string) *model.DBMessage); ok {
		r0 = rf(title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DBMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessagesCount provides a mock function with given fields:
func (_m *MessageRepo) GetMessagesCount() (int, error) {
	ret := _m.Called()

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func() (int, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessagesCountByChannelID provides a mock function with given fields: id
func (_m *MessageRepo) GetMessagesCountByChannelID(id int) (int, error) {
	ret := _m.Called(id)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (int, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) int); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMessageRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageRepo creates a new instance of MessageRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageRepo(t mockConstructorTestingTNewMessageRepo) *MessageRepo {
	mock := &MessageRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
