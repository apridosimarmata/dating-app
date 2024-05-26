package domain

import (
	"dating-app/domain/activity"
	"dating-app/domain/auth"
	"dating-app/domain/feed"
	"dating-app/domain/subscription"
	"dating-app/domain/user"
	"dating-app/infrastructure/utils"
)

type Repositories struct {
	FeedRepository         feed.FeedRepository
	FeedCacheRepository    feed.FeedCacheRepository
	ActivityRepository     activity.ActivityRepository
	UserRepository         user.UserRepository
	UserCacheRepository    user.UserCacheRepository
	SubscriptionRepository subscription.SubscriptionRepository
}

type Usecases struct {
	AuthUsecase         auth.AuthUsecase
	FeedUsecase         feed.FeedUsecase
	UserUsecase         user.UserUsecase
	SubscriptionUsecase subscription.SubscriptionUsecase
}

type Utils struct {
	HashUtils utils.HashUtils
}
