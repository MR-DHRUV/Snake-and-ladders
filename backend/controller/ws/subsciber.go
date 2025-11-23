package ws

import (
    "context"
	"encoding/json"

	"github.com/MR-DHRUV/snake_and_ladders/db"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/MR-DHRUV/snake_and_ladders/transport"
	"github.com/MR-DHRUV/snake_and_ladders/constants"
	gameHandler "github.com/MR-DHRUV/snake_and_ladders/handler/game"
)

func StartGameEventSubscriber() {
    pubsub := db.RedisClient.PSubscribe(context.Background(), constants.RedisChannelPrefix+"*"+constants.RedisChannelSuffix)

	for msg := range pubsub.Channel() {
		gameId := msg.Channel[len(constants.RedisChannelPrefix) : len(msg.Channel)-len(constants.RedisChannelSuffix)]
		evt := model.GameEvent{}
		json.Unmarshal([]byte(msg.Payload), &evt)

		// broadcast to clients of the game
		go gameHandler.HandleBroadcast(&gameHandler.Context{
			ConnectionManager: transport.GetConnectionManager(),
			GameId: gameId,
			UserId: "redis_subscriber",
			UserName: "redis_subscriber",
		}, &evt)
	}
}
