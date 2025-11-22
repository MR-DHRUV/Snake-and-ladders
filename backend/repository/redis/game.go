package redis

import (
	"encoding/json"
	"errors"

	"github.com/MR-DHRUV/snake_and_ladders/db"
	"github.com/MR-DHRUV/snake_and_ladders/model"
)

func GetGameById(gameId string) (*model.Game, error) { 
	cachedGame, err := db.GetCache(gameId)
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

	return db.SetCache(game.Id, string(gameData))
}