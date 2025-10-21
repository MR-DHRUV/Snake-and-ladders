package db

import (
	"context"
	"github.com/MR-DHRUV/snake_and_ladders/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var GameCollection *mongo.Collection

func InitMongo() {
	uri := config.MongoURI
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		utils.GetLogger().Error("Failed to connect to MongoDB: %v", err)
	}
	
	Client = client
	UserCollection = client.Database("snake_and_ladders").Collection("users")
	GameCollection = client.Database("snake_and_ladders").Collection("games")
}
