package feed

import (
	"context"
	"database/sql"
	"dating-app/domain/feed"

	"gorm.io/gorm"
)

type feedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) feed.FeedRepository {
	return &feedRepository{
		db: db,
	}
}

func (repository *feedRepository) GetFeedProfileUids(ctx context.Context, userUid string, excludedUids []string, limit int) (res []string, err error) {
	query := `
	SELECT 
		uid
	FROM users
	WHERE gender != (SELECT gender FROM users WHERE uid = ?) AND uid NOT IN ?
	`

	err = repository.db.WithContext(ctx).Raw(query, userUid, excludedUids).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repository *feedRepository) GetProfileByUids(ctx context.Context, uids []string) (res []feed.FeedProfile, err error) {
	query := `
	SELECT 
		u.uid, 
		u.name, 
		u.gender, 
		u.date_of_birth, 
		u.profile_pic_url, 
		CASE 
			WHEN s.user_uid IS NOT NULL THEN true
			ELSE false
		END AS is_subscriber
	FROM users u
	LEFT JOIN
		(SELECT user_uid FROM subscriptions WHERE expired_at IS NOT NULL AND expired_at > NOW()) AS s
	ON u.uid = s.user_uid
	WHERE u.uid IN ?
	`

	err = repository.db.WithContext(ctx).Raw(query, uids).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
