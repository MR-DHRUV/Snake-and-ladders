package services

import (
	"encoding/json"
	"errors"

	"github.com/MR-DHRUV/snake_and_ladders/constants"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	redis "github.com/MR-DHRUV/snake_and_ladders/repository/redis"
)

func PublishGameEvent(gameId string, event *model.GameEvent) error {
	channel := constants.RedisChannelPrefix + gameId + constants.RedisChannelSuffix
	eventData, err := json.Marshal(event)
	if err != nil {
		return errors.New("failed to marshal game event")
	}

	err = redis.PublishMessage(channel, eventData)
	if err != nil {
		return errors.New("failed to publish game event to Redis")
	}

	return nil
}
