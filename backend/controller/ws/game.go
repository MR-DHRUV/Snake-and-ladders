package ws

import (
	"net/http"

	"github.com/MR-DHRUV/snake_and_ladders/constants"
	gameHandler "github.com/MR-DHRUV/snake_and_ladders/handler/game"
	"github.com/MR-DHRUV/snake_and_ladders/repository"
	"github.com/MR-DHRUV/snake_and_ladders/transport"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
	"github.com/gorilla/websocket"
)

type GameRequest struct {
	Action string `json:"action"`
	Message *string `json:"message,omitempty"`
}

var upgrader = transport.Upgrader
var connectionManager = transport.GetConnectionManager()

func GameWebSocketConnectionController(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie(constants.AuthToken)
    if err != nil {
        http.Error(w, "Authorization cookie missing", http.StatusUnauthorized)
        return
    }

    // Verify JWT
    _, userId, err := utils.VerifyJWT(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

	gameId := r.URL.Query().Get("gameId")
	if gameId == "" || userId == "" {
		http.Error(w, "Missing gameId or userId in URL", http.StatusBadRequest)
		return
	}

	// Upgrade the connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.GetLogger().Error("Failed to upgrade connection: %v", err)
		return
	}

	// Register this connection
	connectionManager.RegisterConnection(gameId, userId, conn)

	// Handle incoming messages in a loop
	go handleMessages(conn, gameId, userId)
}

func handleMessages(conn *websocket.Conn, gameId, userId string) {
	defer func() {
		connectionManager.RemoveConnection(gameId, userId)
		utils.GetLogger().Info("Closing WebSocket connection for GameID: %s, UserID: %s", gameId, userId)
	}()

	for {
		var msg GameRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			utils.GetLogger().Error("Error reading WebSocket message: %v", err)
			break
		}

		utils.GetLogger().Info("Received message: %+v", msg)

		user, err := repository.GetUserById(userId)
		if err != nil {
			user.Name = userId
		}

		ctx := &gameHandler.Context{
			Conn:              conn,
			GameId:            gameId,
			UserId:            userId,
			ConnectionManager: connectionManager,
			Message:           msg.Message,
			UserName:          user.Name,
		}

		switch msg.Action {
		case constants.ActionTypes.JoinGame:
			gameHandler.HandleJoinGame(ctx)
		case constants.ActionTypes.StartGame:
			gameHandler.HandleStartGame(ctx)
		case constants.ActionTypes.NextTurn:
			gameHandler.HandleNextTurn(ctx)
		case constants.ActionTypes.ChatMessage:
			gameHandler.HandleChat(ctx)
		case constants.ActionTypes.RestartGame:
			gameHandler.HandleRestartGame(ctx)
		default:
			resp := utils.ErrorResponse(http.StatusBadRequest, "Unknown action: "+msg.Action)
			conn.WriteJSON(resp)
		}
	}
}
