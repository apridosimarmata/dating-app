package user

import (
	"context"
	"dating-app/domain/user"
	"dating-app/infrastructure"
)

type userCacheRepository struct {
	cache infrastructure.Cache
}

func NewUserCacheRepository(cache infrastructure.Cache) user.UserCacheRepository {
	return &userCacheRepository{
		cache: cache,
	}
}

func (repository *userCacheRepository) SetUserCacheDetails(ctx context.Context, userUid string, user user.UserCacheDetails) (err error) {
	return repository.cache.SetMObj(ctx, map[string]interface{}{userUid: user}, 77)

}

func (repository *userCacheRepository) GetUserCacheDetails(ctx context.Context, userUid string) (res *user.UserCacheDetails, err error) {
	err = repository.cache.GetObj(ctx, userUid, res)
	return res, nil
}
