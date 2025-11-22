package services

import (
	"errors"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/MR-DHRUV/snake_and_ladders/repository"
	redis "github.com/MR-DHRUV/snake_and_ladders/repository/redis"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
)

func CreateNewGame(
	creatorId string,
	maxPlayers,
	maxWinners,
	diceCount int) (*model.Game, error) {

	creator, err := repository.GetUserById(creatorId)
	if err != nil {
		return nil, err
	}

	game := model.NewGame(
		creator,
		maxPlayers,
		maxWinners,
		diceCount,
	)

	updatedGame, err := repository.CreateGame(game)
	if err != nil {
		return nil, err
	}

	game.Id = updatedGame.Id // Id is mongo generated Id

	// Save the game to redis
	err = redis.SetGameById(game)
	if err != nil {
		utils.GetLogger().Error("Failed to cache game in Redis: %v", err)
		return nil, err
	}

	return updatedGame, nil
}

func CreateGameFromLastGame(game_id string) (*string, error) {

	lastGame, err := redis.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	game := model.NewGameFromLastGame(lastGame)
	updatedGame, err := repository.CreateGame(game)
	if err != nil {
		return nil, err
	}

	game.Id = updatedGame.Id
	err = redis.SetGameById(game)
	if err != nil {
		return nil, err
	}

	return &game.Id, nil
}

func JoinGame(user_id, game_id string) (*model.Game, error) {

	user, err := repository.GetUserById(user_id)
	if err != nil {
		return nil, err
	}

	game, err := redis.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	ok, err := game.AddUser(user)
	if err != nil || !ok {
		return nil, err
	}

	err = redis.SetGameById(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func RemoveUser(user_id, game_id string) (*model.Game, error) {
	game, err := redis.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	ok := game.RemoveUser(user_id)
	if !ok {
		return nil, nil
	}

	err = redis.SetGameById(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func StartGame(user_id, game_id string) (*model.Game, error) {
	game, err := redis.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	err = game.StartGame(user_id)
	if err != nil {
		return nil, err
	}

	err = redis.SetGameById(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func NextTurn(user_id, game_id string) (*model.Game, bool, error) {

	game, err := redis.GetGameById(game_id)
	if err != nil {
		return nil, false, err
	}

	ok, err := game.NextTurn(user_id)
	isGameOver := game.IsGameOver()

	if err != nil || !ok {
		return nil, isGameOver, err
	}

	if isGameOver {
		repository.UpdateGame(game) // mark status in db
	}

	err = redis.SetGameById(game)
	if err != nil {
		return nil, isGameOver, err
	}

	return game, isGameOver, nil
}

func GetGames(user_id string, page, limit int) (*model.GetGamesResponse, error) {

	if limit > config.MaxEntriedPerPage {
		return nil, errors.New("limit exceeds maximum allowed value i.e " + string(config.MaxEntriedPerPage))
	}

	return repository.GetGames(user_id, page, limit)
}

// func GetGameById(gameId string) (*model.Game, error) {
// 	game, err := repository.GetGameById(gameId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return game, nil
// }
