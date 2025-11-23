package db

import (
	"context"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

var RedisClient *redis.Client

func InitRedis() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     config.RedisAddr,
        Password: config.RedisPassword,
        DB:       config.RedisDB,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
    })

	RedisClient = rdb
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
    if err != nil {
		utils.GetLogger().Error("Failed to connect to Redis: %v", err)
	} else {
		utils.GetLogger().Info("Connected to Redis")
	}
}
