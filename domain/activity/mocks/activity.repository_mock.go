// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ActivityRepository is an autogenerated mock type for the ActivityRepository type
type ActivityRepository struct {
	mock.Mock
}

// GetUserActivitiesByDate provides a mock function with given fields: ctx, userUid, date
func (_m *ActivityRepository) GetUserActivitiesByDate(ctx context.Context, userUid string, date string) (map[string]string, error) {
	ret := _m.Called(ctx, userUid, date)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) map[string]string); ok {
		r0 = rf(ctx, userUid, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userUid, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserActivityCountByDate provides a mock function with given fields: ctx, userUid, date
func (_m *ActivityRepository) GetUserActivityCountByDate(ctx context.Context, userUid string, date string) (int, error) {
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

// UpdateUserActivitiesByDate provides a mock function with given fields: ctx, userUid, updatedActivities, like, date
func (_m *ActivityRepository) UpdateUserActivitiesByDate(ctx context.Context, userUid string, updatedActivities map[string]string, like map[string]interface{}, date string) error {
	ret := _m.Called(ctx, userUid, updatedActivities, like, date)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string, map[string]interface{}, string) error); ok {
		r0 = rf(ctx, userUid, updatedActivities, like, date)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewActivityRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewActivityRepository creates a new instance of ActivityRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewActivityRepository(t mockConstructorTestingTNewActivityRepository) *ActivityRepository {
	mock := &ActivityRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}