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

func GetCache(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()

	result, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	} else if err != nil {
		return "", err // Some other error occurred
	}

	return result, nil
}

func SetCache(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()

	return RedisClient.Set(ctx, key, value, 0).Err()
}

func DeleteCache(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()
	return RedisClient.Del(ctx, key).Err()
}
