package cache

import (
	"auth/models"
	"github.com/go-redis/redis/v8"
)

type PostCache interface{
	Set(key string, value * models.UserSession)
	Get(key string) (*models.UserSession, error)
	Del(key string) error
	Scan(cursor uint64, match string, count int64) *redis.ScanCmd
}
