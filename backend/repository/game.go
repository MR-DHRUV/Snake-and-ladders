package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/db"
	"github.com/MR-DHRUV/snake_and_ladders/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateGame(game *model.Game) (*model.Game, error) {

	newGame := bson.M{
		"total_cells":  game.TotalCells,
		"max_players":  game.MaxPlayers,
		"max_winners":  game.MaxWinners,
		"dice_count":   game.DiceCount,
		"game_board":   game.GameBoard,
		"status":       game.Status,
		"current_turn": game.CurrentTurn,
		"players":      game.Players,
		"winners":      game.Winners,
		"dice_manager": game.DiceManager,
		"creator":      game.Creator,
		"last_turn":    game.LastTurn,
		"date":         time.Now().UTC().Format(time.RFC3339),
		"messages":     game.Messages,
	}

	// Insert the game into the collection
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	res, err := db.GameCollection.InsertOne(ctx, newGame)
	if err != nil {
		return nil, err
	}

	// Get the inserted game by its generated _id
	var createdGame model.Game
	err = db.GameCollection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&createdGame)
	if err != nil {
		return nil, err
	}

	return &createdGame, nil
}

func UpdateGame(game *model.Game) (*model.Game, error) {
	if game.Id == "" {
		return nil, errors.New("game ID cannot be empty")
	}

	objectID, err := primitive.ObjectIDFromHex(game.Id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	newGame := bson.M{
		"total_cells":  game.TotalCells,
		"max_players":  game.MaxPlayers,
		"max_winners":  game.MaxWinners,
		"dice_count":   game.DiceCount,
		"game_board":   game.GameBoard,
		"status":       game.Status,
		"current_turn": game.CurrentTurn,
		"players":      game.Players,
		"winners":      game.Winners,
		"dice_manager": game.DiceManager,
		"creator":      game.Creator,
		"last_turn":    game.LastTurn,
		"messages":     game.Messages,
	}

	// Prepare the update query
	update := bson.M{"$set": newGame}

	// Update the game document in the collection
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	var updatedGame model.Game
	err = db.GameCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedGame)
	if err != nil {
		return nil, err
	}

	return &updatedGame, nil
}

func GetGameById(gameId string) (*model.Game, error) {

	objectID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	var game model.Game
	filter := bson.M{"_id": objectID}
	err = db.GameCollection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func GetGames(user_id string, page, limit int) (*model.GetGamesResponse, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("page and limit must be greater than 0")
	}

	skip := (page - 1) * limit
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	filter := bson.M{
		"players._id": user_id,
	}

	selectFields := bson.M{
		"_id":     1,
		"date":    1,
		"creator": 1,
		"status":  1,
		"players": 1,
		"winners": 1,
	}

	options := options.Find().SetProjection(selectFields).SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{"date", -1}})
	
	cursor, err := db.GameCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var games []*model.GameResponse
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}

	// Get the total count of games matching the filter
	totalCount, err := db.GameCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &model.GetGamesResponse{
		Games: games,
		Total: totalCount,
	}, nil
}