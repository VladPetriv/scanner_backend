// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/VladPetriv/scanner_backend/internal/model"

// MessageService is an autogenerated mock type for the MessageService type
type MessageService struct {
	mock.Mock
}

// GetFullMessageByMessageID provides a mock function with given fields: ID
func (_m *MessageService) GetFullMessageByMessageID(ID int) (*model.FullMessage, error) {
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
func (_m *MessageService) GetFullMessages(page int) ([]model.FullMessage, error) {
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
func (_m *MessageService) GetFullMessagesByChannelID(ID int, limit int, page int) ([]model.FullMessage, error) {
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
func (_m *MessageService) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
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

// GetMessage provides a mock function with given fields: messagelID
func (_m *MessageService) GetMessage(messagelID int) (*model.Message, error) {
	ret := _m.Called(messagelID)

	var r0 *model.Message
	if rf, ok := ret.Get(0).(func(int) *model.Message); ok {
		r0 = rf(messagelID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(messagelID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessageByName provides a mock function with given fields: name
func (_m *MessageService) GetMessageByName(name string) (*model.Message, error) {
	ret := _m.Called(name)

	var r0 *model.Message
	if rf, ok := ret.Get(0).(func(string) *model.Message); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessages provides a mock function with given fields:
func (_m *MessageService) GetMessages() ([]model.Message, error) {
	ret := _m.Called()

	var r0 []model.Message
	if rf, ok := ret.Get(0).(func() []model.Message); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
