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

// UpdateGameWithLock updates the game state using optimistic locking (Redis WATCH).
// The modifier function receives the current game state and should return the modified game state (or nil to abort) and a boolean indicating if the update should proceed.
func UpdateGameWithLock(gameId string, modifier func(*model.Game) (*model.Game, error)) (*model.Game, error) {
	ctx := context.Background()
	var updatedGame *model.Game

	// Retry loop for optimistic locking
	err := db.RedisClient.Watch(ctx, func(tx *redis.Tx) error {
		// 1. Get the current game state
		gameStr, err := tx.Get(ctx, gameId).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		var game *model.Game
		if gameStr != "" {
			game = &model.Game{}
			err = json.Unmarshal([]byte(gameStr), game)
			if err != nil {
				return err
			}
		} else {
			return errors.New("game not found")
		}

		// 2. Apply the modification
		modifiedGame, err := modifier(game)
		if err != nil {
			return err
		}
		if modifiedGame == nil {
			return nil // Abort update without error
		}

		gameData, err := json.Marshal(modifiedGame)
		if err != nil {
			return err
		}

		// 3. Save the modified game state in a transaction
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, gameId, string(gameData), 0)
			return nil
		})

		if err == nil {
			updatedGame = modifiedGame
		}

		return err
	}, gameId)

	if err != nil {
		return nil, err
	}

	return updatedGame, nil
}
