// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	feed "dating-app/domain/feed"

	mock "github.com/stretchr/testify/mock"
)

// FeedCacheRepository is an autogenerated mock type for the FeedCacheRepository type
type FeedCacheRepository struct {
	mock.Mock
}

// DeleteSubscriberFeedProfileUids provides a mock function with given fields: ctx, userUid
func (_m *FeedCacheRepository) DeleteSubscriberFeedProfileUids(ctx context.Context, userUid string) error {
	ret := _m.Called(ctx, userUid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userUid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProfileByUids provides a mock function with given fields: ctx, uids
func (_m *FeedCacheRepository) GetProfileByUids(ctx context.Context, uids []string) ([]feed.FeedProfile, error) {
	ret := _m.Called(ctx, uids)

	var r0 []feed.FeedProfile
	if rf, ok := ret.Get(0).(func(context.Context, []string) []feed.FeedProfile); ok {
		r0 = rf(ctx, uids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]feed.FeedProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, uids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriberFeedProfileUids provides a mock function with given fields: ctx, userUid
func (_m *FeedCacheRepository) GetSubscriberFeedProfileUids(ctx context.Context, userUid string) ([]string, error) {
	ret := _m.Called(ctx, userUid)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, userUid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userUid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetProfiles provides a mock function with given fields: ctx, profiles
func (_m *FeedCacheRepository) SetProfiles(ctx context.Context, profiles map[string]feed.FeedProfile) error {
	ret := _m.Called(ctx, profiles)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]feed.FeedProfile) error); ok {
		r0 = rf(ctx, profiles)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetSubscriberFeedProfileUids provides a mock function with given fields: ctx, userUid, profileUids
func (_m *FeedCacheRepository) SetSubscriberFeedProfileUids(ctx context.Context, userUid string, profileUids []string) error {
	ret := _m.Called(ctx, userUid, profileUids)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) error); ok {
		r0 = rf(ctx, userUid, profileUids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewFeedCacheRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewFeedCacheRepository creates a new instance of FeedCacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFeedCacheRepository(t mockConstructorTestingTNewFeedCacheRepository) *FeedCacheRepository {
	mock := &FeedCacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}