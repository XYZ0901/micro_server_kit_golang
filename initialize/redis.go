package initialize

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func redisInit() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Cfg.RedisConfig.Addr,
		Password: Cfg.RedisConfig.PassWord,
		DB:       0,
		PoolSize: 100,
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
