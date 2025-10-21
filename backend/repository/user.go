package repository

import (
	"context"
	"time"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/db"
	"github.com/MR-DHRUV/snake_and_ladders/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveUser(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	newUser := bson.M{
		"name":    user.Name,
		"email":   user.Email,
		"picture": user.Picture,
		"date":    time.Now().UTC().Format(time.RFC3339),
	}

	result, err := db.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	var savedUser model.User
	err = db.UserCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&savedUser)
	if err != nil {
		return nil, err
	}

	return &savedUser, nil
}

func GetUserById(userId string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	var user model.User
	filter := bson.M{"_id": objectID}
	err = db.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	filter := bson.M{"email": email}

	var user model.User
	err := db.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
