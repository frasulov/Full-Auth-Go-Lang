package cache

import (
	"auth/models/redisModels"
	"github.com/go-redis/redis/v8"
)

type PostCache interface {
	Set(key string, value *redisModels.UserSession)
	Get(key string) (*redisModels.UserSession, error)
	Del(key string) error
	Scan(cursor uint64, match string, count int64) *redis.ScanCmd
}
