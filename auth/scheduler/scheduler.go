package scheduler

import (
	"auth/config"
	cache "auth/redis"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

var scheduler *cron.Cron

func init(){
	scheduler = cron.New()
}

func DeleteExpiredSessionFromRedis() {
	scheduler.AddFunc(config.Configuration.SessionDeleteScheduler, func() {
		redis := cache.NewRedisCache(fmt.Sprintf("%s:%v", config.Configuration.Redis.Host, config.Configuration.Redis.Port), config.Configuration.Redis.Db, time.Duration(config.Configuration.Redis.Expires))
		iter := redis.Scan(0, "accessToken-of-*", 0).Iterator()
		for iter.Next(context.Background()) {
			splitToken := strings.Split(iter.Val(), "accessToken-of-")
			userSession ,err := redis.Get(fmt.Sprintf("accessToken-of-%v", splitToken[1]))
			if err != nil {
				panic(err)
			}
			if userSession.ExpiredAt <= time.Now().Unix() {
				fmt.Println("Session expired Deleting")
				redis.Del(fmt.Sprintf("accessToken-of-%v", splitToken[1]))
			}else{
				fmt.Println("Session is active")
			}
		}
		if err := iter.Err(); err != nil {
			panic(err)
		}
	})
}

func StartScheduler(){
	scheduler.Start()
}
