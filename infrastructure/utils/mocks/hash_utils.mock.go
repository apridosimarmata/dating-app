// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// HashUtils is an autogenerated mock type for the HashUtils type
type HashUtils struct {
	mock.Mock
}

// HashPassword provides a mock function with given fields: plainPassword, salt
func (_m *HashUtils) HashPassword(plainPassword string, salt string) string {
	ret := _m.Called(plainPassword, salt)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(plainPassword, salt)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewHashUtils interface {
	mock.TestingT
	Cleanup(func())
}

// NewHashUtils creates a new instance of HashUtils. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHashUtils(t mockConstructorTestingTNewHashUtils) *HashUtils {
	mock := &HashUtils{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
