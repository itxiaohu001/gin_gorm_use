package database

import (
	"test/config"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})

	// 可以在这里添加一个 Ping 操作来验证连接
	_, err := RedisClient.Ping(RedisClient.Context()).Result()
	return err
}
