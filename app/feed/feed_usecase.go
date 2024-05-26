package feed

import (
	"context"
	"dating-app/domain"
	"dating-app/domain/activity"
	"dating-app/domain/common/response"
	"dating-app/domain/feed"
	"dating-app/domain/user"
	"dating-app/infrastructure"
	"errors"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
)

const (
	activityLockKey = "activity:%s"
)

type feedUsecase struct {
	feedRepository      feed.FeedRepository
	feedCacheRepository feed.FeedCacheRepository
	activityRepository  activity.ActivityRepository
	userRepository      user.UserRepository
	userCacheRepository user.UserCacheRepository
	mutexProvider       infrastructure.MutexProvider
}

func NewFeedUsecase(repositories domain.Repositories, mutexProvider infrastructure.MutexProvider) feed.FeedUsecase {
	return &feedUsecase{
		feedRepository:      repositories.FeedRepository,
		feedCacheRepository: repositories.FeedCacheRepository,
		activityRepository:  repositories.ActivityRepository,
		userRepository:      repositories.UserRepository,
		userCacheRepository: repositories.UserCacheRepository,
		mutexProvider:       mutexProvider,
	}
}

func (usecase *feedUsecase) LikeProfile(ctx context.Context, userUid string, profileUid string) (res *response.Response[interface{}]) {
	var activityLock *redsync.Mutex
	var err error
	todayDate := time.Now().Format("02-01-2006")

	if activityLock, err = usecase.mutexProvider.Acquire(fmt.Sprintf(activityLockKey, userUid)); err != nil {
		infrastructure.Log("got error on usecase.mutexProvider.Acquire() - LikeProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_ANOTHER_PROCESS_IS_RUNNING,
		}
	}

	activities, err := usecase.activityRepository.GetUserActivitiesByDate(ctx, userUid, todayDate)
	if err != nil {
		infrastructure.Log("got error on usecase.activityRepository.GetUserActivitiesByDate() - LikeProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	activities[profileUid] = activity.LIKE_ACTIVITY
	like := make(map[string]interface{}, 1)
	like[profileUid] = nil

	err = usecase.activityRepository.UpdateUserActivitiesByDate(ctx, userUid, activities, like, todayDate)
	if err != nil {
		infrastructure.Log("got error on usecase.activityRepository.UpdateUserActivitiesByDate() - LikeProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if err := usecase.mutexProvider.Release(activityLock); err != nil {
		infrastructure.Log("got error on usecase.mutexProvider.Release() - LikeProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	return &response.Response[interface{}]{
		Error: false,
	}
}

func (usecase *feedUsecase) SkipProfile(ctx context.Context, userUid string, profileUid string) (res *response.Response[interface{}]) {
	var activityLock *redsync.Mutex
	var err error
	todayDate := time.Now().Format("02-01-2006")

	if activityLock, err = usecase.mutexProvider.Acquire(fmt.Sprintf(activityLockKey, userUid)); err != nil {
		infrastructure.Log("got error on usecase.mutexProvider.Acquire() - SkipProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_ANOTHER_PROCESS_IS_RUNNING,
		}
	}

	activities, err := usecase.activityRepository.GetUserActivitiesByDate(ctx, userUid, todayDate)
	if err != nil {
		infrastructure.Log("got error on usecase.activityRepository.GetUserActivitiesByDate() - SkipProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	activities[profileUid] = activity.SKIP_ACTIVITY
	like := make(map[string]interface{}, 1)
	like[profileUid] = nil

	err = usecase.activityRepository.UpdateUserActivitiesByDate(ctx, userUid, activities, nil, todayDate)
	if err != nil {
		infrastructure.Log("got error on usecase.activityRepository.UpdateUserActivitiesByDate() - SkipProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if err := usecase.mutexProvider.Release(activityLock); err != nil {
		infrastructure.Log("got error on usecase.mutexProvider.Release() - SkipProfile")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	return &response.Response[interface{}]{
		Error: false,
	}
}

func (usecase *feedUsecase) GetProfileFeeds(ctx context.Context, userUid string) (res *response.Response[[]feed.FeedProfile]) {
	profiles := []feed.FeedProfile{}

	user, err := usecase.userCacheRepository.GetUserCacheDetails(ctx, userUid)
	if err != nil {
		user, err = usecase.userRepository.GetUserMiniDetailsByUserUid(ctx, userUid)
		if err != nil {
			infrastructure.Log("got error on usecase.userRepository.GetUserMiniDetailsByUserUid() - GetProfileFeeds")
			return &response.Response[[]feed.FeedProfile]{
				Error:   true,
				Message: &response.ERROR_INTERNAL_SERVER_ERROR,
			}
		}
	}

	if user.IsSubscriber {
		profiles, err = usecase.getSubscriberFeedProfiles(ctx, userUid)
		if err != nil {
			infrastructure.Log("got error on usecase.getSubscriberFeedProfiles() - GetProfileFeeds")
			return &response.Response[[]feed.FeedProfile]{
				Error:   true,
				Message: &response.ERROR_INTERNAL_SERVER_ERROR,
			}
		}
	} else {
		profiles, err = usecase.getNonSubscriberFeedProfiles(ctx, userUid)
		if err != nil {
			infrastructure.Log("got error on usecase.getSubscriberFeedProfiles() - GetProfileFeeds")
			errMsg := err.Error()
			return &response.Response[[]feed.FeedProfile]{
				Error:   true,
				Message: &errMsg,
			}
		}
	}

	return &response.Response[[]feed.FeedProfile]{
		Error: false,
		Data:  &profiles,
	}
}

func (usecase *feedUsecase) getNonSubscriberFeedProfiles(ctx context.Context, userUid string) (res []feed.FeedProfile, err error) {
	todayDate := time.Now().Format("02-01-2006")

	userTodayActivities, err := usecase.activityRepository.GetUserActivitiesByDate(ctx, userUid, todayDate)
	if err != nil {
		infrastructure.Log("got error on usecase.activityRepository.GetUserActivitiesByDate() - getNonSubscriberFeedProfiles")
		return nil, err
	}

	if len(userTodayActivities) >= activity.REGULAR_USER_DAILY_LIMIT {
		return nil, errors.New(response.ERROR_REACHED_LIMIT)
	}

	swipedProfileUids := []string{}
	for profileUid, swipeType := range userTodayActivities {
		if swipeType == activity.SKIP_ACTIVITY {
			swipedProfileUids = append(swipedProfileUids, profileUid)
		}
	}
	swipedProfileUids = append(swipedProfileUids, userUid)

	profileUids, err := usecase.feedRepository.GetFeedProfileUids(ctx, userUid, swipedProfileUids, activity.REGULAR_USER_DAILY_LIMIT-len(userTodayActivities))
	if err != nil {
		infrastructure.Log("got error on usecase.feedRepository.GetFeedProfileUids() - getNonSubscriberFeedProfiles")
		return nil, err
	}

	return usecase.getFeedProfileByProfileUids(ctx, profileUids)
}

func (usecase *feedUsecase) getSubscriberFeedProfiles(ctx context.Context, userUid string) (res []feed.FeedProfile, err error) {
	var profileUidsToReturn []string
	todayDate := time.Now().Format("02-01-2006")

	userTodayActivities, err := usecase.activityRepository.GetUserActivitiesByDate(ctx, userUid, todayDate)
	if err != nil {
		infrastructure.Log("got error on usecase.activityRepository.GetUserActivitiesByDate() - getSubscriberFeedProfiles")
		return nil, err
	}

	swipedProfileUids := []string{}
	for profileUid, swipeType := range userTodayActivities {
		if swipeType == activity.SKIP_ACTIVITY {
			swipedProfileUids = append(swipedProfileUids, profileUid)
		}
	}
	swipedProfileUids = append(swipedProfileUids, userUid)

	profileUids, err := usecase.feedCacheRepository.GetSubscriberFeedProfileUids(ctx, userUid)
	if err != nil {
		profileUids, err = usecase.feedRepository.GetFeedProfileUids(ctx, userUid, swipedProfileUids, 100)
		if err != nil {
			infrastructure.Log("got error on usecase.feedRepository.GetFeedProfileUids() - getSubscriberFeedProfiles")
			return nil, err
		}
	}

	// there are profile uids left after this call, cache them
	if len(profileUids) > 10 {
		profileUidsToReturn = profileUids[0:10]
		err := usecase.feedCacheRepository.SetSubscriberFeedProfileUids(ctx, userUid, profileUids[10:])
		if err != nil {
			infrastructure.Log("got error on usecase.feedCacheRepository.SetSubscriberFeedProfileUids() - getSubscriberFeedProfiles")
			return nil, err
		}
	} else {
		// delete key since there are no more profile uids being stored
		profileUidsToReturn = profileUids
		err := usecase.feedCacheRepository.DeleteSubscriberFeedProfileUids(ctx, userUid)
		if err != nil {
			infrastructure.Log("got error on usecase.feedCacheRepository.DeleteSubscriberFeedProfileUids() - getSubscriberFeedProfiles")
			return nil, err
		}
	}

	return usecase.getFeedProfileByProfileUids(ctx, profileUidsToReturn)
}

func (usecase *feedUsecase) getFeedProfileByProfileUids(ctx context.Context, profileUids []string) (res []feed.FeedProfile, err error) {
	res, err = usecase.feedCacheRepository.GetProfileByUids(ctx, profileUids)
	if err != nil {
		infrastructure.Log("got error on usecase.feedCacheRepository.GetProfileByUids() - getFeedProfileByProfileUids")
		return nil, err
	}

	if len(profileUids) != len(res) {
		getProfileFromDBUids := []string{}

		cachedProfileUidsMap := make(map[string]struct{}, len(profileUids))
		for _, cachedProfile := range res {
			cachedProfileUidsMap[cachedProfile.UserUid] = struct{}{}
		}

		for _, uid := range profileUids {
			if _, cached := cachedProfileUidsMap[uid]; !cached {
				getProfileFromDBUids = append(getProfileFromDBUids, uid)
			}
		}

		_profiles, err := usecase.feedRepository.GetProfileByUids(ctx, getProfileFromDBUids)
		if err != nil {
			infrastructure.Log("got error on usecase.feedRepository.GetProfileByUids() - getFeedProfileByProfileUids")
			return nil, err
		}

		cacheProfileToSet := make(map[string]feed.FeedProfile, len(_profiles))
		for _, profile := range _profiles {
			cacheProfileToSet[profile.UserUid] = profile
		}
		err = usecase.feedCacheRepository.SetProfiles(ctx, cacheProfileToSet)
		if err != nil {
			infrastructure.Log("got error on usecase.feedCacheRepository.SetProfiless() - getFeedProfileByProfileUids")
		}

		res = append(res, _profiles...)
	}

	return res, nil
}
