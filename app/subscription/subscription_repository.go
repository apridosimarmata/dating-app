package subscription

import (
	"context"
	"database/sql"
	"dating-app/domain/subscription"

	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) subscription.SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (repository *subscriptionRepository) InsertSubscription(ctx context.Context, subscription subscription.Subscription) (err error) {
	err = repository.db.Table("subscriptions").Create(&subscription).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *subscriptionRepository) GetSubscriptionByUid(ctx context.Context, uid string) (res *subscription.Subscription, err error) {
	query := `
		SELECT 
			uid, 
			user_uid, 
			created_at, 
			expired_at, 
			package_id, 
			paid_at, 
			payment_amount, 
			payment_status 
		FROM subscriptions
		WHERE uid = ?
	`

	err = repository.db.WithContext(ctx).Raw(query, uid).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (repository *subscriptionRepository) UpdateSubscription(ctx context.Context, subscription subscription.Subscription) (err error) {
	err = repository.db.WithContext(ctx).Table("subscriptions").UpdateColumns(subscription).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *subscriptionRepository) GetPackageById(ctx context.Context, packageId int) (res *subscription.Package, err error) {
	query := `
		select id, duration, price from packages where id = ?
	`

	err = repository.db.WithContext(ctx).Raw(query, packageId).Scan(&res).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
