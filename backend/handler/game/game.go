package game

import (
	"net/http"
	"strings"

	"github.com/MR-DHRUV/snake_and_ladders/constants"
	"github.com/MR-DHRUV/snake_and_ladders/model"
	"github.com/MR-DHRUV/snake_and_ladders/services"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
)

func HandleJoinGame(ctx *Context) {
	game, err := services.JoinGame(ctx.UserId, ctx.GameId)
	if err != nil {
		ctx.Conn.WriteJSON(utils.ErrorResponse(http.StatusInternalServerError, "failed to join game: "+err.Error()))
		return
	}

	// notify all players in the game
	utils.GetLogger().Info("Game state: %v", game.GetGameState())
	utils.GetLogger().Info("Messages: %v", game.Messages)
	broadcastGameState(game.GetGameState(), ctx)
}

func HandleNextTurn(ctx *Context) {
	game, _, err := services.NextTurn(ctx.UserId, ctx.GameId)
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

	broadcastGameState(game.GetGameState(), ctx)
}

func HandleStartGame(ctx *Context) {
	game, err := services.StartGame(ctx.UserId, ctx.GameId)
	if err != nil {
		ctx.Conn.WriteJSON(utils.ErrorResponse(http.StatusInternalServerError, "failed to start game: "+err.Error()))
		return
	}

	broadcastGameState(game.GetGameState(), ctx)
}

func HandleRestartGame(ctx *Context) {
	nextGameId, err := services.CreateGameFromLastGame(ctx.GameId)
	if err != nil {
		ctx.Conn.WriteJSON(utils.ErrorResponse(http.StatusInternalServerError, "failed to restart game: "+err.Error()))
		return
	}

	broadcast(GameResponse{
		Type: constants.ResponseTypes.NextGame,
		Data: utils.NextGameMessage{GameId: *nextGameId},
	}, ctx)
}

func HandleChat(ctx *Context) {
	if ctx.Message == nil || strings.TrimSpace(*ctx.Message) == "" {
		return
	}

	resp := GameResponse{
		Type: constants.ResponseTypes.ChatMessage,
		Data: utils.ChatMessage{
			PlayerName: ctx.UserName,
			PlayerId:   ctx.UserId,
			Message:    *ctx.Message,
		},
	}

	broadcast(resp, ctx)
}

func handleBrokenConnections(ctx *Context) {
	if ctx.UserId == "" {
		return
	}

	game, err := services.RemoveUser(ctx.UserId, ctx.GameId)
	if err == nil {
		broadcastGameState(game.GetGameState(), ctx)
	}
}

// BroadcastGameResponse broadcasts the game state to all connected players
func broadcastGameState(gameResponse model.GameStateChangeResponse, ctx *Context) {
	broadcast(GameResponse{
		Type: constants.ResponseTypes.GameState,
		Data: gameResponse,
	}, ctx)
}

func broadcast(resp GameResponse, ctx *Context) {
	connections := ctx.ConnectionManager.GetGameConnections(ctx.GameId)
	if connections == nil {
		utils.GetLogger().Info("No connections found for game %s", ctx.GameId)
		return
	}

	brokenConnections := make([]string, 0)
	for userId, conn := range connections {
		err := conn.WriteJSON(resp)
		if err != nil {
			utils.GetLogger().Error("Failed to send message to player %s: %v", userId, err)
			brokenConnections = append(brokenConnections, userId)
			ctx.ConnectionManager.RemoveConnection(ctx.GameId, userId)
		}
	}
}
