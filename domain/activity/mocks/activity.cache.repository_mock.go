// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ActivityCacheRepository is an autogenerated mock type for the ActivityCacheRepository type
type ActivityCacheRepository struct {
	mock.Mock
}

// GetUserActivityCountByDate provides a mock function with given fields: ctx, userUid, date
func (_m *ActivityCacheRepository) GetUserActivityCountByDate(ctx context.Context, userUid string, date string) (int, error) {
	ret := _m.Called(ctx, userUid, date)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int); ok {
		r0 = rf(ctx, userUid, date)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userUid, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetUserTodayActivityCount provides a mock function with given fields: ctx, userUid, count
func (_m *ActivityCacheRepository) SetUserTodayActivityCount(ctx context.Context, userUid string, count int) error {
	ret := _m.Called(ctx, userUid, count)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) error); ok {
		r0 = rf(ctx, userUid, count)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewActivityCacheRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewActivityCacheRepository creates a new instance of ActivityCacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewActivityCacheRepository(t mockConstructorTestingTNewActivityCacheRepository) *ActivityCacheRepository {
	mock := &ActivityCacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}