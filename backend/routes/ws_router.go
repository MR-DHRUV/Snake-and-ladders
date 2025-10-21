package routes

import (
	"fmt"
	"net/http"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/utils"

	gameController "github.com/MR-DHRUV/snake_and_ladders/controller/ws"

)

func StartWebSocketServer() {
	http.HandleFunc("/game", gameController.GameWebSocketConnectionController)
	utils.GetLogger().Info("WebSocket server is running on port %d", config.WebSocketPort)

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.WebSocketPort), nil)
	if err != nil {
		utils.GetLogger().Error("Failed to start WebSocket server: %v", err)
	}
}