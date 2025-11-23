package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/db"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/redis/go-redis/v9"
)

func getCache(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()

	result, err := db.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	} else if err != nil {
		return "", err // Some other error occurred
	}

	return result, nil
}

func setCache(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()

	return db.RedisClient.Set(ctx, key, value, 0).Err()
}

func deleteCache(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()
	return db.RedisClient.Del(ctx, key).Err()
}

func GetGameById(gameId string) (*model.Game, error) {
	cachedGame, err := getCache(gameId)
	if err != nil {
		return nil, err
	}

	if cachedGame == "" {
		return nil, errors.New("game not found in cache")
	}

	game := &model.Game{}
	err = json.Unmarshal([]byte(cachedGame), game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func SetGameById(game *model.Game) error {
	gameData, err := json.Marshal(game)
	if err != nil {
		return err
	}

	return setCache(game.Id, string(gameData))
}

func DeleteGameById(gameId string) error {
	return deleteCache(gameId)
}
