package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func NewRedisClient(ctx context.Context, config Config) redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return *rdb
}

type Cache interface {
	SetString(ctx context.Context, key string, obj string, ttlInSec int) (err error)
	GetString(ctx context.Context, key string) (result string, err error)
	Del(ctx context.Context, key string) (err error)
	GetObj(ctx context.Context, key string, dest interface{}) (err error)
	SetMObj(ctx context.Context, payload map[string]interface{}, ttlInSec int) (err error)
	GetMObj(ctx context.Context, params []string) (map[string]interface{}, error)
	Publish(ctx context.Context, channel string, payload interface{}) (err error)
}

type redisCache struct {
	client redis.Client
}

func NewCache(redisClient redis.Client) Cache {
	return &redisCache{
		client: redisClient,
	}
}

func (cache *redisCache) GetMObj(ctx context.Context, params []string) (map[string]interface{}, error) {
	cmd := cache.client.MGet(params...)
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for pos, key := range params {
		if res[pos] != nil {
			result[key] = res[pos]
		}
	}

	return result, nil
}

func (cache *redisCache) GetObj(ctx context.Context, key string, dest interface{}) (err error) {
	result, err := cache.client.Get(key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		return err
	}

	return nil
}

func (cache *redisCache) SetMObj(ctx context.Context, payload map[string]interface{}, ttlInSec int) (err error) {
	pipe := cache.client.Pipeline()
	for key, val := range payload {
		b, err := json.Marshal(val)
		if err != nil {
			continue
		}
		pipe.Set(key, b, time.Duration(ttlInSec)*time.Second)
	}
	_, err = pipe.Exec()
	if err != nil {
		return err
	}
	return nil

}

func (cache *redisCache) Publish(ctx context.Context, channel string, payload interface{}) (err error) {
	payloadInString, err := json.Marshal(payload)
	if err != nil {
		return
	}

	_, err = cache.client.Publish(channel, payloadInString).Result()
	return err
}

func (cache *redisCache) SetString(ctx context.Context, key string, obj string, ttlInSec int) (err error) {
	return cache.client.Set(key, obj, time.Second*time.Duration(ttlInSec)).Err()
}

func (cache *redisCache) GetString(ctx context.Context, key string) (result string, err error) {
	res, err := cache.client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (cache *redisCache) Del(ctx context.Context, key string) (err error) {

	return cache.client.Del(key).Err()
}
