package transport

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for simplicity (but you can restrict this for production)
		return true
	},
}

var (
	connectionManager *ConnectionManager
	once              sync.Once
)

// GetConnectionManager: singleton pattern to get the ConnectionManager instance
func GetConnectionManager() *ConnectionManager {
	once.Do(func() {
		connectionManager = NewConnectionManager()
	})
	return connectionManager
}