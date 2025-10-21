package game

import (
	"github.com/MR-DHRUV/snake_and_ladders/transport"
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

type GameResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

