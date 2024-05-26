package subscription

import (
	"context"
	"dating-app/domain/common/response"
)

const (
	PAYMENT_VAT = .11
)

type Subscription struct {
	UID           string  `json:"uid" gorm:"column=uid"`
	UserUid       string  `json:"user_uid" gorm:"column=user_uid"`
	CreatedAt     string  `json:"created_at" gorm:"created_at"`
	ExpiredAt     *string `json:"expired_at" gorm:"expired_at"`
	PackageId     int     `json:"package_id" gorm:"column=package_id"`
	PaidAt        *string `json:"paid_at" gorm:"column=paid_at"`
	PaymentAmount int     `json:"payment_amount" gorm:"column=payment_amount"`
	PaymentStatus string  `json:"payment_status" gorm:"column=payment_status"`
}

type SubscriptionPaymentCallback struct {
	UID    string `json:"uid"`
	Status string `json:"status"`
}

type CreateSubscriptionRequest struct {
	UserUid   string `json:"user_uid"`
	PackageId int    `json:"package_id"`
}

type Package struct {
	ID       int `json:"id" gorm:"column=id"`
	Duration int `json:"duration" gorm:"column=duration"`
	Price    int `json:"price" gorm:"column=price"`
}

type SubscriptionRepository interface {
	InsertSubscription(ctx context.Context, subscription Subscription) (err error)
	GetSubscriptionByUid(ctx context.Context, uid string) (res *Subscription, err error)
	UpdateSubscription(ctx context.Context, subscription Subscription) (err error)

	GetPackageById(ctx context.Context, packageId int) (res *Package, err error)
}

type SubscriptionUsecase interface {
	Subscribe(ctx context.Context, request CreateSubscriptionRequest) (res *response.Response[interface{}])
	HandleSubscriptionPaymentCallback(ctx context.Context, request SubscriptionPaymentCallback) (res *response.Response[interface{}])
}
