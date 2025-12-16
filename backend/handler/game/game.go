package game

import (
	"net/http"
	"strings"

	"github.com/MR-DHRUV/snake_and_ladders/constants"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/MR-DHRUV/snake_and_ladders/services"
	"github.com/MR-DHRUV/snake_and_ladders/transport"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
	"github.com/gorilla/websocket"
)

type Context struct {
	ConnectionManager *transport.ConnectionManager
	Conn              *websocket.Conn
	Message           *string
	GameId            string
	UserId            string
	UserName          string
}

func HandleJoinGame(ctx *Context) {
	game, err := services.JoinGame(ctx.UserId, ctx.GameId)
	if err != nil {
		ctx.Conn.WriteJSON(utils.ErrorResponse(http.StatusInternalServerError, "failed to join game: "+err.Error()))
		return
	}

	// notify all players in the game
	utils.GetLogger().Info("Game state: %v", game.GetGameState())
	utils.GetLogger().Info("Messages: %v", game.Messages)
	publishGameStateBroadcastEvent(ctx)
}

func HandleNextTurn(ctx *Context) {
	_, _, err := services.NextTurn(ctx.UserId, ctx.GameId)
	if err != nil {
		if strings.Contains(err.Error(), "it's not your turn") {
			resp := utils.ErrorResponse(http.StatusForbidden, "It's not your turn")
			ctx.Conn.WriteJSON(resp)
		} else {
			resp := utils.ErrorResponse(http.StatusInternalServerError, "Failed to process game turn: "+err.Error())
			ctx.Conn.WriteJSON(resp)
		}
		return
	}

	publishGameStateBroadcastEvent(ctx)
}

func HandleStartGame(ctx *Context) {
	_, err := services.StartGame(ctx.UserId, ctx.GameId)
	if err != nil {
		ctx.Conn.WriteJSON(utils.ErrorResponse(http.StatusInternalServerError, "failed to start game: "+err.Error()))
		return
	}

	publishGameStateBroadcastEvent(ctx)
}

func HandleRestartGame(ctx *Context) {
	nextGameId, err := services.CreateGameFromLastGame(ctx.GameId)
	if err != nil {
		ctx.Conn.WriteJSON(utils.ErrorResponse(http.StatusInternalServerError, "failed to restart game: "+err.Error()))
		return
	}

	event := model.GameEvent{
		Type: constants.ResponseTypes.NextGame,
		Data: utils.NextGameMessage{GameId: *nextGameId},
	}

	services.PublishGameEvent(ctx.GameId, &event)
}

func HandleChat(ctx *Context) {
	if ctx.Message == nil || strings.TrimSpace(*ctx.Message) == "" {
		return
	}

	event := model.GameEvent{
		Type: constants.ResponseTypes.ChatMessage,
		Data: utils.ChatMessage{
			PlayerName: ctx.UserName,
			PlayerId:   ctx.UserId,
			Message:    *ctx.Message,
		},
	}

	services.PublishGameEvent(ctx.GameId, &event)
}

func HandleBroadcast(ctx *Context, event *model.GameEvent) {
	if( event == nil || ctx == nil ) {
		return
	}

	if event.Type == constants.ResponseTypes.GameState {
		game, err := services.GetGameById(ctx.GameId)
		if err != nil {
			utils.GetLogger().Error("Failed to get game for broadcasting game state: %v", err)
			return
		}

		event.Data = game.GetGameState()
	}

	connections := ctx.ConnectionManager.GetGameConnections(ctx.GameId)
	if connections == nil {
		utils.GetLogger().Info("No connections found for game %s", ctx.GameId)
		return
	}

	brokenConnections := make([]string, 0)
	for userId, conn := range connections {
		err := conn.WriteJSON(event)
		if err != nil {
			utils.GetLogger().Error("Failed to send message to player %s: %v", userId, err)
			brokenConnections = append(brokenConnections, userId)
			ctx.ConnectionManager.RemoveConnection(ctx.GameId, userId)
		}
	}
}

func handleBrokenConnections(ctx *Context) {
	if ctx.UserId == "" {
		return
	}

	_, err := services.RemoveUser(ctx.UserId, ctx.GameId)
	if err == nil {
		publishGameStateBroadcastEvent(ctx)
	}
}

func publishGameStateBroadcastEvent(ctx *Context) {
	services.PublishGameEvent(ctx.GameId, &model.GameEvent{
		Type: constants.ResponseTypes.GameState,
		Data: nil,
	})
}
