package user

import (
	"context"
	"database/sql"
	"dating-app/domain/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repository *userRepository) GetUserDetailsByEmail(ctx context.Context, email string) (res *user.User, err error) {
	query := `
	SELECT
		uid,
		name,
		email,
		gender,
		profile_pic_url,
		password,
		dob,
		created_at,
		updated_at
	FROM
		users where email = ?
	`

	err = repository.db.WithContext(ctx).Raw(query, email).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repository *userRepository) GetUserCountByUid(ctx context.Context, uid string) (res int, err error) {
	query := `
	select count(*) from users where uid = ?

	`

	err = repository.db.WithContext(ctx).Raw(query, uid).Scan(&res).Error
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (repository *userRepository) GetUserCountByEmail(ctx context.Context, email string) (res int, err error) {
	query := `
		select count(*) from users where email = ?
	`

	err = repository.db.WithContext(ctx).Raw(query, email).Scan(&res).Error
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (repository *userRepository) InsertUser(ctx context.Context, user user.User) (err error) {
	err = repository.db.Table("users").Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) GetUserMiniDetailsByUserUid(ctx context.Context, userUid string) (res *user.UserCacheDetails, err error) {
	query := `
	SELECT 
		u.uid, 
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

	err = repository.db.WithContext(ctx).Raw(query, userUid).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
