package cache

import (
	"auth/models/redisModels"
	"context"
	json2 "encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, expires time.Duration) PostCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: expires,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	client := cache.getClient()
	return client.Scan(context.Background(), cursor, match, count)
}

func (cache *redisCache) Set(key string, value *redisModels.UserSession) {
	client := cache.getClient()
	json, err := json2.Marshal(value)
	if err != nil {
		panic(err.Error())
	}

	client.Set(context.Background(), key, json, cache.expires*time.Second)
}

func (cache *redisCache) Del(key string) error {
	client := cache.getClient()
	_, err := client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (cache *redisCache) Get(key string) (*redisModels.UserSession, error) {
	client := cache.getClient()
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return &redisModels.UserSession{}, err
	}
	var post redisModels.UserSession
	err = json2.Unmarshal([]byte(val), &post)
	if err != nil {
		return &redisModels.UserSession{}, err
	}
	return &post, nil
}
