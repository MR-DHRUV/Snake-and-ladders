package services

import (
	"errors"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/MR-DHRUV/snake_and_ladders/repository"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
)

func CreateNewGame(
	creatorId string,
	maxPlayers,
	maxWinners,
	diceCount int) (*model.Game, error) {

	utils.GetLogger().Info("Creating new game with parameters: "+
		"CreatorId: %s, "+
		"MaxPlayers: %d, MaxWinners: %d, DiceCount: %d",
		creatorId,
		maxPlayers, maxWinners, diceCount)

	creator, err := repository.GetUserById(creatorId)
	if err != nil {
		return nil, err
	}

	utils.GetLogger().Info("Fetched creator user: %s", creator.Name)

	game := model.NewGame(
		creator,
		maxPlayers,
		maxWinners,
		diceCount,
	)

	utils.GetLogger().Info("Creating new game with ID: %s", game.Id)

	updatedGame, err := repository.CreateGame(game)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil
}

func CreateGameFromLastGame(game_id string) (*string, error) {

	lastGame, err := repository.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	game := model.NewGameFromLastGame(lastGame)
	game, err = repository.CreateGame(game)
	if err != nil {
		return nil, err
	}

	return &game.Id, nil;
}

func JoinGame(user_id, game_id string) (*model.Game, error) {

	user, err := repository.GetUserById(user_id)
	if err != nil {
		return nil, err
	}

	game, err := repository.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	ok, err := game.AddUser(user)
	if err != nil || !ok {
		return nil, err
	}

	updatedGame, err := repository.UpdateGame(game)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil
}

func RemoveUser(user_id, game_id string) (*model.Game, error) {
	game, err := repository.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	ok := game.RemoveUser(user_id)
	if !ok {
		return nil, nil
	}

	updatedGame, err := repository.UpdateGame(game)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil
}

func StartGame(user_id, game_id string) (*model.Game, error) {
	game, err := repository.GetGameById(game_id)
	if err != nil {
		return nil, err
	}

	err = game.StartGame(user_id)
	if err != nil {
		return nil, err
	}

	updatedGame, err := repository.UpdateGame(game)
	if err != nil {
		return nil, err
	}

	return updatedGame, nil
}

func NextTurn(user_id, game_id string) (*model.Game, bool, error) {

	game, err := repository.GetGameById(game_id)
	if err != nil {
		return nil, false, err
	}

	ok, isGameOver, err := game.NextTurn(user_id)
	if err != nil || !ok {
		return nil, isGameOver, err
	}

	updatedGame, err := repository.UpdateGame(game)
	if err != nil {
		return nil, isGameOver, err
	}

	return updatedGame, isGameOver, nil
}

func GetGames(user_id string, page, limit int) (*model.GetGamesResponse, error) {

	if limit > config.MaxEntriedPerPage {
		return nil, errors.New("limit exceeds maximum allowed value i.e " + string(config.MaxEntriedPerPage))
	}

	return repository.GetGames(user_id, page, limit)
}

func GetGameById(gameId string) (*model.Game, error) {
	game, err := repository.GetGameById(gameId)
	if err != nil {
		return nil, err
	}

	return game, nil
}
