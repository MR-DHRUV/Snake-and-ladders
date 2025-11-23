package redis

import (
	"context"
	"fmt"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/db"
	"github.com/redis/go-redis/v9"
)

func PublishMessage(channel string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()

	return db.RedisClient.Publish(ctx, channel, message).Err()
}

func SubscribeToChannel(channel string) (*redis.PubSub, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.RedisTimeout)
	defer cancel()

	pubsub := db.RedisClient.Subscribe(ctx, channel)

	// Wait for confirmation that subscription is created before returning
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel %s: %v", channel, err)
	}

	return pubsub, nil
}
