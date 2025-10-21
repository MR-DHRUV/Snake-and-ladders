const ActionTypes = {
    JoinGame: "joinGame",
    StartGame: "startGame",
    NextTurn: "nextTurn",
    ChatMessage: "chatMessage",
    RestartGame: "restartGame",
};

const ResponseTypes = {
    GameState: "gameState",
    ChatMessage: "chatMessage",
    NextGame: "nextGame"
}

export {
    ActionTypes,
    ResponseTypes
}