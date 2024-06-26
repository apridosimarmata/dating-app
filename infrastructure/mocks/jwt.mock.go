// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// JwtInterface is an autogenerated mock type for the JwtInterface type
type JwtInterface struct {
	mock.Mock
}

// GenerateAccessToken provides a mock function with given fields: userUid, jwtSecret, isRefreshToken
func (_m *JwtInterface) GenerateAccessToken(userUid string, jwtSecret string, isRefreshToken bool) (*string, error) {
	ret := _m.Called(userUid, jwtSecret, isRefreshToken)

	var r0 *string
	if rf, ok := ret.Get(0).(func(string, string, bool) *string); ok {
		r0 = rf(userUid, jwtSecret, isRefreshToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, bool) error); ok {
		r1 = rf(userUid, jwtSecret, isRefreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateToken provides a mock function with given fields: tokenString, jwtSecret
func (_m *JwtInterface) ValidateToken(tokenString string, jwtSecret string) (*string, int64, error) {
	ret := _m.Called(tokenString, jwtSecret)

	var r0 *string
	if rf, ok := ret.Get(0).(func(string, string) *string); ok {
		r0 = rf(tokenString, jwtSecret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(string, string) int64); ok {
		r1 = rf(tokenString, jwtSecret)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(tokenString, jwtSecret)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewJwtInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewJwtInterface creates a new instance of JwtInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJwtInterface(t mockConstructorTestingTNewJwtInterface) *JwtInterface {
	mock := &JwtInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
