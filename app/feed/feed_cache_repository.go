package feed

import (
	"context"
	"dating-app/domain/feed"
	"dating-app/infrastructure"
	"encoding/json"
)

type feedCacheRepository struct {
	cache infrastructure.Cache
}

func NewFeedCacheRepository(cache infrastructure.Cache) feed.FeedCacheRepository {
	return &feedCacheRepository{
		cache: cache,
	}
}

func (repository *feedCacheRepository) DeleteSubscriberFeedProfileUids(ctx context.Context, userUid string) (err error) {
	return repository.cache.Del(ctx, userUid)
}

func (repository *feedCacheRepository) SetSubscriberFeedProfileUids(ctx context.Context, userUid string, profileUids []string) (err error) {
	return repository.cache.SetMObj(ctx, map[string]interface{}{userUid: profileUids}, 77)
}

func (repository *feedCacheRepository) GetSubscriberFeedProfileUids(ctx context.Context, userUid string) (res []string, err error) {
	err = repository.cache.GetObj(ctx, userUid, res)
	return res, nil
}

func (repository *feedCacheRepository) GetProfileByUids(ctx context.Context, uids []string) (res []feed.FeedProfile, err error) {
	cacheObj, err := repository.cache.GetMObj(ctx, uids)
	if err != nil {
		return nil, err
	}

	for _, item := range cacheObj {
		dest := feed.FeedProfile{}
		err = json.Unmarshal([]byte(item.(string)), dest)
		if err != nil {
			return nil, err
		}

		res = append(res, dest)
	}

	return res, nil

}

func (repository *feedCacheRepository) SetProfiles(ctx context.Context, profiles map[string]feed.FeedProfile) (err error) {
	payload := map[string]interface{}{}

	for uid, profile := range profiles {
		payload[uid] = profile
	}

	return repository.cache.SetMObj(ctx, payload, 77)
}
