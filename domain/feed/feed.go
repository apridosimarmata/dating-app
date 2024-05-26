package feed

import (
	"context"
	"dating-app/domain/common/response"
)

type FeedProfile struct {
	UserUid       string `json:"user_uid" gorm:"column=uid"`
	Name          string `json:"name" gorm:"column=name"`
	Gender        rune   `json:"gender" gorm:"column=gender"`
	DateOfBirth   string `json:"dob" gorm:"column=dob"`
	ProfilePicUrl string `json:"profile_pic_url" gorm:"column=profile_pic_url"`
}

type FeedCacheRepository interface {
	DeleteSubscriberFeedProfileUids(ctx context.Context, userUid string) (err error)
	SetSubscriberFeedProfileUids(ctx context.Context, userUid string, profileUids []string) (err error)
	GetSubscriberFeedProfileUids(ctx context.Context, userUid string) (res []string, err error)
	GetProfileByUids(ctx context.Context, uids []string) (res []FeedProfile, err error)
	SetProfiles(ctx context.Context, profiles map[string]FeedProfile) (err error)
}

type FeedRepository interface {
	GetFeedProfileUids(ctx context.Context, userUid string, excludedUids []string, limit int) (res []string, err error)
	GetProfileByUids(ctx context.Context, uids []string) (res []FeedProfile, err error)
}

type FeedUsecase interface {
	GetProfileFeeds(ctx context.Context, userUid string) (res *response.Response[[]FeedProfile])
	LikeProfile(ctx context.Context, userUid string, profileUid string) (res *response.Response[interface{}])
	SkipProfile(ctx context.Context, userUid string, profileUid string) (res *response.Response[interface{}]) // equals to skip profile
}
