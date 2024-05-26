package subscription

import (
	"context"
	"dating-app/domain"
	"dating-app/domain/common/response"
	"dating-app/domain/subscription"
	"dating-app/domain/user"
	"dating-app/infrastructure"
	"time"

	"github.com/google/uuid"
)

type subscriptionUsecase struct {
	userRepository         user.UserRepository
	subscriptionRepository subscription.SubscriptionRepository
}

func NewSubscriptionUsecase(repositories domain.Repositories) subscription.SubscriptionUsecase {
	return &subscriptionUsecase{
		userRepository:         repositories.UserRepository,
		subscriptionRepository: repositories.SubscriptionRepository,
	}
}

func (usecase *subscriptionUsecase) Subscribe(ctx context.Context, request subscription.CreateSubscriptionRequest) (res *response.Response[interface{}]) {
	userCount, err := usecase.userRepository.GetUserCountByUid(ctx, request.UserUid)
	if err != nil {
		infrastructure.Log("got error on usecase.userRepository.GetUserCountByUid() - Subscribe")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if userCount == 0 {
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_USER_NOT_FOUND,
		}
	}

	_package, err := usecase.subscriptionRepository.GetPackageById(ctx, request.PackageId)
	if err != nil {
		infrastructure.Log("got error on usecase.subscriptionRepository.GetPackageById() - Subscribe")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if _package == nil {
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_PACKAGE_NOT_FOUND,
		}
	}

	subscriptionUid, err := uuid.NewV6()
	if err != nil {
		infrastructure.Log("got error on uuid.NewV6() - RegisterUser")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	_subscription := subscription.Subscription{
		UID:           subscriptionUid.String(),
		UserUid:       request.UserUid,
		PackageId:     request.PackageId,
		PaymentAmount: int(float64(_package.Price) + (float64(_package.Price) * subscription.PAYMENT_VAT)),
		PaymentStatus: "CREATED",
	}

	err = usecase.subscriptionRepository.InsertSubscription(ctx, _subscription)
	if err != nil {
		infrastructure.Log("got error on usecase.subscriptionRepository.InsertSubscription() - Subscribe")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	return &response.Response[interface{}]{
		Error: false,
	}
}
func (usecase *subscriptionUsecase) HandleSubscriptionPaymentCallback(ctx context.Context, request subscription.SubscriptionPaymentCallback) (res *response.Response[interface{}]) {
	_subscription, err := usecase.subscriptionRepository.GetSubscriptionByUid(ctx, request.UID)
	if err != nil {
		infrastructure.Log("got error on usecase.subscriptionRepository.GetSubscriptionByUid() - Subscribe")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if _subscription == nil {
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_SUBSCRIPTION_NOT_FOUND,
		}
	}

	_subscription.PaymentStatus = request.Status
	_package, err := usecase.subscriptionRepository.GetPackageById(ctx, _subscription.PackageId)
	if err != nil {
		infrastructure.Log("got error on usecase.subscriptionRepository.GetPackageById() - Subscribe")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	if request.Status == "PAID" {
		expiredAt := int(time.Now().Unix()) + (24 * 30 * _package.Duration * int(time.Hour))
		expiredAtString := time.Unix(int64(expiredAt), 0).Format(time.RFC3339)
		_subscription.ExpiredAt = &expiredAtString
	}

	err = usecase.subscriptionRepository.UpdateSubscription(ctx, *_subscription)
	if err != nil {
		infrastructure.Log("got error on usecase.subscriptionRepository.GetSubscriptionByUid() - Subscribe")
		return &response.Response[interface{}]{
			Error:   true,
			Message: &response.ERROR_INTERNAL_SERVER_ERROR,
		}
	}

	return &response.Response[interface{}]{
		Error: false,
	}
}
