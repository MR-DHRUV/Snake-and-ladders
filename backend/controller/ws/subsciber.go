package ws

import (
	"sync"
	"context"
	"encoding/json"

	"github.com/MR-DHRUV/snake_and_ladders/constants"
	"github.com/MR-DHRUV/snake_and_ladders/db"
	gameHandler "github.com/MR-DHRUV/snake_and_ladders/handler/game"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/MR-DHRUV/snake_and_ladders/transport"
)

type gameChannel struct {
	EventChan chan *model.GameEvent
}

var channelMap sync.Map // map[gameId]*gameChannel

func getOrCreateWorker(gameId string) *gameChannel {
	if w, ok := channelMap.Load(gameId); ok {
		return w.(*gameChannel)
	}

	// create new worker
	worker := &gameChannel{
		EventChan: make(chan *model.GameEvent, 100),
	}

	channelMap.Store(gameId, worker)

	// launch worker goroutine
	go startWorker(gameId, worker)

	return worker
}

func startWorker(gameId string, worker *gameChannel) {
	ctx := &gameHandler.Context{
		ConnectionManager: transport.GetConnectionManager(),
		GameId:            gameId,
		UserId:            "redis_worker",
		UserName:          "redis_worker",
	}

	for evt := range worker.EventChan {
		gameHandler.HandleBroadcast(ctx, evt)
	}
}

func StartGameEventSubscriber() {
	pubsub := db.RedisClient.PSubscribe(context.Background(), constants.RedisChannelPrefix+"*"+constants.RedisChannelSuffix)

	for msg := range pubsub.Channel() {
		gameId := msg.Channel[len(constants.RedisChannelPrefix) : len(msg.Channel)-len(constants.RedisChannelSuffix)]

		evt := model.GameEvent{}
		json.Unmarshal([]byte(msg.Payload), &evt)

		// Push into worker queue for sequential processing
		worker := getOrCreateWorker(gameId)
		worker.EventChan <- &evt
	}
}
