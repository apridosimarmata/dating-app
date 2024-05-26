package activity

import "context"

const (
	SKIP_ACTIVITY = "SKIP"
	LIKE_ACTIVITY = "LIKE"

	REGULAR_USER_DAILY_LIMIT = 10
)

type Activities struct {
	UserUid         string                            `json:"user_uid"`
	DailyActivities map[string]map[string]string      `json:"daily_activity"`
	Pass            map[string]map[string]interface{} `json:"pass"` // e.g. 22/03/2024 -> { "uid1": null, "uid2": null, ...}
}

type Likes struct {
	UserUid string                 `json:"user_uid"`
	Likes   map[string]interface{} `json:"swipes"`
}

type ActivityRepository interface {
	GetUserActivityCountByDate(ctx context.Context, userUid string, date string) (res int, err error)
	GetUserActivitiesByDate(ctx context.Context, userUid string, date string) (res map[string]string, err error)
	UpdateUserActivitiesByDate(ctx context.Context, userUid string, updatedActivities map[string]string, like map[string]interface{}, date string) (err error)
}
