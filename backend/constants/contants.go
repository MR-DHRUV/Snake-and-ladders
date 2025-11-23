package constants

var ActionTypes = struct {
	JoinGame    string
	NextTurn    string
	StartGame   string
	ChatMessage string
	RestartGame string
}{
	JoinGame:    "joinGame",
	StartGame:   "startGame",
	NextTurn:    "nextTurn",
	ChatMessage: "chatMessage",
	RestartGame: "restartGame",
}

var ResponseTypes = struct {
	NextGame   string
	GameState  string
	ChatMessage string
}{
	NextGame:   "nextGame",
	GameState:  "gameState",
	ChatMessage: "chatMessage",
}

var RedisChannelPrefix = "game-"
var RedisChannelSuffix = "-events"
var AuthToken = "token"
