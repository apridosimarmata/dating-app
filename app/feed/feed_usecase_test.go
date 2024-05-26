package feed

import (
	"context"
	"dating-app/domain"
	"dating-app/domain/activity"
	activityMocks "dating-app/domain/activity/mocks"
	"dating-app/domain/common/response"
	"dating-app/domain/feed"
	feedMocks "dating-app/domain/feed/mocks"
	"dating-app/domain/user"

	userMocks "dating-app/domain/user/mocks"
	infraMocks "dating-app/infrastructure/mocks"
	"testing"

	"github.com/go-redsync/redsync/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type feedTestParams struct {
	feedRepository          *feedMocks.FeedRepository
	feedCacheRepository     *feedMocks.FeedCacheRepository
	activityRepository      *activityMocks.ActivityRepository
	activityCacheRepository *activityMocks.ActivityCacheRepository
	userRepository          *userMocks.UserRepository
	userCacheRepository     *userMocks.UserCacheRepository
	mutexProvider           *infraMocks.MutexProvider
	context                 context.Context
	feedUsecase             feed.FeedUsecase
}

func setupFeedTestParams(t *testing.T) feedTestParams {
	feedRepository := feedMocks.NewFeedRepository(t)
	feedCacheRepository := feedMocks.NewFeedCacheRepository(t)
	activityRepository := activityMocks.NewActivityRepository(t)
	userRepository := userMocks.NewUserRepository(t)
	userCacheRepository := userMocks.NewUserCacheRepository(t)
	mutexProvider := infraMocks.NewMutexProvider(t)

	feedUsecase := NewFeedUsecase(domain.Repositories{
		FeedRepository:      feedRepository,
		FeedCacheRepository: feedCacheRepository,
		ActivityRepository:  activityRepository,
		UserRepository:      userRepository,
		UserCacheRepository: userCacheRepository,
	},
		mutexProvider)

	return feedTestParams{
		feedRepository:      feedRepository,
		feedCacheRepository: feedCacheRepository,
		activityRepository:  activityRepository,
		userRepository:      userRepository,
		userCacheRepository: userCacheRepository,
		context:             context.Background(),
		feedUsecase:         feedUsecase,
		mutexProvider:       mutexProvider,
	}
}

func returnError() error {
	return assert.AnError
}

func returnNil() error {
	return assert.AnError
}

func Test_LikeProfile(t *testing.T) {
	params := setupFeedTestParams(t)

	t.Run("return error on failed acquire lock", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.LikeProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetUserActivitiesByDate error", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.LikeProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetUserActivitiesByDate error", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.activityRepository.On("UpdateUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()

		res := params.feedUsecase.LikeProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on Release error", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.activityRepository.On("UpdateUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		params.mutexProvider.On("Release", mock.Anything).Return(assert.AnError).Once()

		res := params.feedUsecase.LikeProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return success", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.activityRepository.On("UpdateUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		params.mutexProvider.On("Release", mock.Anything).Return(nil).Once()

		res := params.feedUsecase.LikeProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, false, res.Error)
	})
}

func Test_SkipProfile(t *testing.T) {
	params := setupFeedTestParams(t)

	t.Run("return error on failed acquire lock", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.SkipProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetUserActivitiesByDate error", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.SkipProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetUserActivitiesByDate error", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.activityRepository.On("UpdateUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()

		res := params.feedUsecase.SkipProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on Release error", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.activityRepository.On("UpdateUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		params.mutexProvider.On("Release", mock.Anything).Return(assert.AnError).Once()

		res := params.feedUsecase.SkipProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return success", func(t *testing.T) {
		params.mutexProvider.On("Acquire", mock.Anything).Return(&redsync.Mutex{}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.activityRepository.On("UpdateUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		params.mutexProvider.On("Release", mock.Anything).Return(nil).Once()

		res := params.feedUsecase.SkipProfile(params.context, mock.Anything, mock.Anything)
		require.Equal(t, false, res.Error)
	})
}

func Test_GetProfileFeeds_Subscriber(t *testing.T) {
	params := setupFeedTestParams(t)

	t.Run("return error on GetUserMiniDetailsByUserUid error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetUserActivitiesByDate error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetFeedProfileUids error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.feedCacheRepository.On("GetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on SetSubscriberFeedProfileUids error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{
			"uid20": activity.SKIP_ACTIVITY,
		}, nil).Once()
		params.feedCacheRepository.On("GetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{
			"uid1",
			"uid2",
			"uid3",
			"uid4",
			"uid5",
			"uid6",
			"uid7",
			"uid8",
			"uid9",
			"uid10",
			"uid11",
			"uid12",
		}, nil).Once()
		params.feedCacheRepository.On("SetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on DeleteSubscriberFeedProfileUids error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.feedCacheRepository.On("GetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{"uid1"}, nil).Once()
		params.feedCacheRepository.On("DeleteSubscriberFeedProfileUids", mock.Anything, mock.Anything).Return(assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetProfileByUids [cache] error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.feedCacheRepository.On("GetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{"uid1"}, nil).Once()
		params.feedCacheRepository.On("DeleteSubscriberFeedProfileUids", mock.Anything, mock.Anything).Return(nil).Once()
		params.feedCacheRepository.On("GetProfileByUids", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on GetProfileByUids [DB] error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.feedCacheRepository.On("GetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{"uid1", "uid2"}, nil).Once()
		params.feedCacheRepository.On("DeleteSubscriberFeedProfileUids", mock.Anything, mock.Anything).Return(nil).Once()
		params.feedCacheRepository.On("GetProfileByUids", mock.Anything, mock.Anything).Return([]feed.FeedProfile{
			{
				UserUid: "uid1",
			},
		}, nil).Once()
		params.feedRepository.On("GetProfileByUids", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return success - despite of SetProfiles error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.feedCacheRepository.On("GetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{"uid1", "uid2"}, nil).Once()
		params.feedCacheRepository.On("DeleteSubscriberFeedProfileUids", mock.Anything, mock.Anything).Return(nil).Once()
		params.feedCacheRepository.On("GetProfileByUids", mock.Anything, mock.Anything).Return([]feed.FeedProfile{
			{
				UserUid: "uid1",
			},
		}, nil).Once()
		params.feedRepository.On("GetProfileByUids", mock.Anything, mock.Anything).Return([]feed.FeedProfile{
			{
				UserUid: "uid2",
			},
		}, nil).Once()
		params.feedCacheRepository.On("SetProfiles", mock.Anything, mock.Anything).Return(assert.AnError)

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, false, res.Error)
	})

	t.Run("return success", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(nil, assert.AnError).Once()
		params.userRepository.On("GetUserMiniDetailsByUserUid", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: true,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{}, nil).Once()
		params.feedCacheRepository.On("GetSubscriberFeedProfileUids", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{"uid1", "uid2"}, nil).Once()
		params.feedCacheRepository.On("DeleteSubscriberFeedProfileUids", mock.Anything, mock.Anything).Return(nil).Once()
		params.feedCacheRepository.On("GetProfileByUids", mock.Anything, mock.Anything).Return([]feed.FeedProfile{
			{
				UserUid: "uid1",
			},
			{
				UserUid: "uid2",
			},
		}, nil).Once()
		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, false, res.Error)
	})

}

func Test_GetProfileFeeds_NonSubscriber(t *testing.T) {
	params := setupFeedTestParams(t)

	t.Run("return error on GetUserActivitiesByDate error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: false,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return error on reached daily activity limit", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: false,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{
			"uid1":  activity.SKIP_ACTIVITY,
			"uid2":  activity.SKIP_ACTIVITY,
			"uid3":  activity.LIKE_ACTIVITY,
			"uid4":  activity.SKIP_ACTIVITY,
			"uid5":  activity.LIKE_ACTIVITY,
			"uid6":  activity.SKIP_ACTIVITY,
			"uid7":  activity.LIKE_ACTIVITY,
			"uid8":  activity.LIKE_ACTIVITY,
			"uid9":  activity.SKIP_ACTIVITY,
			"uid10": activity.SKIP_ACTIVITY,
		}, nil).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
		require.Equal(t, response.ERROR_REACHED_LIMIT, *res.Message)

	})

	t.Run("return error on GetFeedProfileUids error", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: false,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{
			"uid1": activity.SKIP_ACTIVITY,
			"uid2": activity.SKIP_ACTIVITY,
			"uid3": activity.LIKE_ACTIVITY,
			"uid4": activity.SKIP_ACTIVITY,
			"uid5": activity.LIKE_ACTIVITY,
		}, nil).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, true, res.Error)
	})

	t.Run("return success", func(t *testing.T) {
		params.userCacheRepository.On("GetUserCacheDetails", mock.Anything).Return(&user.UserCacheDetails{
			IsSubscriber: false,
		}, nil).Once()
		params.activityRepository.On("GetUserActivitiesByDate", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{
			"uid1": activity.SKIP_ACTIVITY,
			"uid2": activity.SKIP_ACTIVITY,
			"uid3": activity.LIKE_ACTIVITY,
			"uid4": activity.SKIP_ACTIVITY,
			"uid5": activity.LIKE_ACTIVITY,
		}, nil).Once()
		params.feedRepository.On("GetFeedProfileUids", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{"uid6", "uid7"}, nil).Once()
		params.feedCacheRepository.On("GetProfileByUids", mock.Anything, mock.Anything).Return([]feed.FeedProfile{
			{
				UserUid: "uid6",
			},
			{
				UserUid: "uid7",
			},
		}, nil).Once()

		res := params.feedUsecase.GetProfileFeeds(params.context, mock.Anything)
		require.Equal(t, false, res.Error)
	})

}
