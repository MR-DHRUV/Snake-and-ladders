package transport

import (
    "sync"

    "github.com/MR-DHRUV/snake_and_ladders/utils"
    "github.com/gorilla/websocket"
)

// ConnectionManager manages WebSocket connections for all games
type ConnectionManager struct {
    // Map of gameId -> Map of userId -> WebSocket connection
    connections map[string]map[string]*websocket.Conn
    mu          sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
    return &ConnectionManager{
        connections: make(map[string]map[string]*websocket.Conn),
    }
}

// RegisterConnection registers a new WebSocket connection for a player in a game
func (cm *ConnectionManager) RegisterConnection(gameId string, userId string, conn *websocket.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    if _, ok := cm.connections[gameId]; !ok {
        cm.connections[gameId] = make(map[string]*websocket.Conn)
    }

    // Close existing connection if any
    if existingConn, exists := cm.connections[gameId][userId]; exists {
        existingConn.Close()
    }

    cm.connections[gameId][userId] = conn
}

// RemoveConnection removes a WebSocket connection
func (cm *ConnectionManager) RemoveConnection(gameId string, userId string) {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    if gameConnections, ok := cm.connections[gameId]; ok {
        if conn, exists := gameConnections[userId]; exists {
            conn.Close()
            delete(gameConnections, userId)
            utils.GetLogger().Info("Removed WebSocket connection for game %s, user %s", gameId, userId)
        }
        
        // If no more connections for this game, remove the game entry
        if len(gameConnections) == 0 {
            delete(cm.connections, gameId)
        }
    }
}

func (cm *ConnectionManager) RemoveGameConnections(gameId string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if gameConnections, ok := cm.connections[gameId]; ok {
		for userId, conn := range gameConnections {
			conn.Close()
			utils.GetLogger().Info("Removed WebSocket connection for game %s, user %s", gameId, userId)
		}
		delete(cm.connections, gameId)
		utils.GetLogger().Info("Removed all WebSocket connections for game %s", gameId)
	}
}


// GetConnection retrieves a WebSocket connection
func (cm *ConnectionManager) GetConnection(gameId string, userId string) (*websocket.Conn, bool) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    if gameConnections, ok := cm.connections[gameId]; ok {
        conn, exists := gameConnections[userId]
        return conn, exists
    }
    return nil, false
}

// GetGameConnections retrieves all connections for a game
func (cm *ConnectionManager) GetGameConnections(gameId string) map[string]*websocket.Conn {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    if gameConnections, ok := cm.connections[gameId]; ok {
        // Make a copy to avoid race conditions
        connsCopy := make(map[string]*websocket.Conn, len(gameConnections))
        for userId, conn := range gameConnections {
            connsCopy[userId] = conn
        }
        return connsCopy
    }
    return nil
}