package model

import (
	"errors"
	"fmt"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
)

const (
	Created = iota
	InProgress
	Finished
	Abandoned
)

type Player struct {
	Id       string `json:"_id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Picture  string `json:"picture" bson:"picture"`
	Position int    `json:"position" bson:"position"`
}

type LastTurn struct {
	PlayerName       string `json:"player_name" bson:"player_name"`
	DiceRolls        []int  `json:"dice_rolls" bson:"dice_rolls"`
	DiceRoll         int    `json:"-" bson:"-"`
	PreviousPosition int    `json:"previous_position" bson:"previous_position"`
	NextPosition     int    `json:"next_position" bson:"next_position"`
}

type GameStateChangeResponse struct {
	Id          string              `json:"_id"`
	Creator     *User               `json:"creator"`
	CurrentTurn string              `json:"current_turn"`
	Players     []*Player           `json:"players"`
	Winners     []*Player           `json:"winners"`
	Status      int                 `json:"status"`
	GameBoard   *Board              `json:"game_board"`
	LastTurn    *LastTurn           `json:"last_turn,omitempty"`
	Messages    []utils.ChatMessage `json:"messages"`
	DiceCount   int                 `json:"dice_count"`
}

type Game struct {
	Id          string              `bson:"_id,omitempty" json:"_id"`
	Date        string              `bson:"date" json:"date"`
	Creator     *User               `bson:"creator" json:"creator"`
	TotalCells  int                 `bson:"total_cells" json:"total_cells"`
	MaxPlayers  int                 `bson:"max_players" json:"max_players"`
	MaxWinners  int                 `bson:"max_winners" json:"max_winners"`
	DiceCount   int                 `bson:"dice_count" json:"dice_count"`
	GameBoard   *Board              `bson:"game_board" json:"game_board"`
	Status      int                 `bson:"status" json:"status"`
	CurrentTurn int                 `bson:"current_turn" json:"current_turn"`
	Players     []*Player           `bson:"players" json:"players"`
	Winners     []*Player           `bson:"winners" json:"winners"`
	DiceManager *DiceManager        `bson:"dice_manager" json:"dice_manager"`
	LastTurn    *LastTurn           `bson:"last_turn" json:"last_turn,omitempty"`
	Messages    []utils.ChatMessage `bson:"-" json:"-"`
}

type GameResponse struct {
	Id      string    `json:"_id" bson:"_id"`
	Date    string    `json:"date" bson:"date"`
	Creator *User     `json:"creator" bson:"creator"`
	Status  int       `json:"status" bson:"status"`
	Players []*Player `json:"players" bson:"players"`
	Winners []*Player `json:"winners" bson:"winners"`
}

type GetGamesResponse struct {
	Games []*GameResponse `json:"games"`
	Total int64           `json:"total"`
}

func (g *Game) AddUser(user *User) (bool, error) {
	existingUser := false
	for _, u := range g.Players {
		if u.Id == user.Id {
			existingUser = true
			break
		}
	}

	// maybe someone is reconnecting
	if existingUser {
		g.Messages = []utils.ChatMessage{
			getSystemMessage(fmt.Sprintf("%s rejoined the game", user.Name)),
		}
		return true, nil
	}

	// Check if the maximum players limit is reached
	if len(g.Players) == g.MaxPlayers {
		return false, errors.New("maximum players limit reached")
	}

	if g.Status != Created {
		return false, errors.New("game is started or finished")
	}

	g.Players = append(g.Players, &Player{
		Id:       user.Id,
		Name:     user.Name,
		Picture:  user.Picture,
		Position: 0, // Initialize player position to 0
	})

	g.Messages = []utils.ChatMessage{
		getSystemMessage(fmt.Sprintf("%s joined the game", user.Name)),
	}

	return true, nil
}

// TODO: check if it is causing memory leak
func (g *Game) RemoveUser(user_id string) bool {
	for i, u := range g.Players {
		if u.Id == user_id {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			break
		}
	}

	if len(g.Players) < 2 {
		g.Status = Abandoned
	}

	return true
}

func (g *Game) GetNextTurn() string {
	return g.Players[g.CurrentTurn].Id
}

func (g *Game) IsGameOver() bool {
	if g.Status == Abandoned {
		return true
	}

	if g.Status == Finished {
		return true
	}

	if g.Status == InProgress && (len(g.Winners) == g.MaxWinners || len(g.Players) < 2) {
		g.Status = Finished
		return true
	}

	return false
}

func (g *Game) StartGame(user_id string) error {
	if g.Status != Created {
		return errors.New("game is already started or finished")
	}

	if len(g.Players) < 2 {
		return errors.New("not enough players to start the game")
	}

	if g.Creator.Id != user_id {
		return errors.New("only the creator can start the game")
	}

	g.Status = InProgress
	g.Messages = []utils.ChatMessage{
		getSystemMessage("Game started!"),
	}

	return nil
}

func (g *Game) NextTurn(_id string) (bool, error) {
	if len(g.Players) == 0 {
		return false, errors.New("no users in the turn queue")
	}

	if g.Status != InProgress {
		return false, errors.New("game is not started")
	}

	if g.GetNextTurn() != _id {
		return false, errors.New("it's not your turn")
	}

	currUser := g.Players[g.CurrentTurn]
	currPosition := currUser.Position

	diceRoll, rolls := g.DiceManager.RollAll()
	nextPosition := g.GameBoard.GetCell(currPosition + diceRoll).GetNextPosition()
	currUser.Position = nextPosition

	// Another turn if all dice rolled 6
	if diceRoll != g.DiceCount*6 {
		g.CurrentTurn++
		if g.CurrentTurn == len(g.Players) {
			g.CurrentTurn = 0
		}
	}

	g.LastTurn = &LastTurn{
		PlayerName:       currUser.Name,
		DiceRolls:        rolls,
		DiceRoll:         diceRoll,
		PreviousPosition: currPosition,
		NextPosition:     nextPosition,
	}

	if nextPosition == g.TotalCells {
		g.Winners = append(g.Winners, currUser)
	}

	g.Messages = []utils.ChatMessage{
		getSystemMessage(fmt.Sprintf("%s rolled %d and moved from %d to %d", currUser.Name, diceRoll, currPosition, nextPosition)),
	}

	if g.IsGameOver() {
		g.Messages = append(g.Messages, getSystemMessage("Game over!"))
	} else {
		g.Messages = append(g.Messages, getSystemMessage(fmt.Sprintf("%s's turn", g.Players[g.CurrentTurn].Name)))
	}

	return true, nil
}

func (g *Game) GetGameState() GameStateChangeResponse {

	messages := g.Messages
	if messages == nil {
		messages = []utils.ChatMessage{}
	}

	return GameStateChangeResponse{
		Id:          g.Id,
		CurrentTurn: g.GetNextTurn(),
		Players:     g.Players,
		Winners:     g.Winners,
		Status:      g.Status,
		GameBoard:   g.GameBoard,
		LastTurn:    g.LastTurn,
		Messages:    messages,
		Creator:     g.Creator,
		DiceCount:   g.DiceCount,
	}
}

func NewGame(
	creator *User,
	maxPlayers,
	maxWinners,
	diceCount int) *Game {
	board := NewBoard(config.BoardCells)
	diceManager := NewDiceManager(diceCount, 6)
	winners := make([]*Player, 0, maxWinners)

	game := &Game{
		TotalCells:  config.BoardCells,
		MaxPlayers:  maxPlayers,
		MaxWinners:  maxWinners,
		DiceCount:   diceCount,
		GameBoard:   board,
		Status:      Created,
		CurrentTurn: 0,
		DiceManager: diceManager,
		Creator:     creator,
		Winners:     winners,
	}

	game.AddUser(creator)
	game.Messages = []utils.ChatMessage{
		getSystemMessage("Game created!"),
	}

	return game
}

func NewGameFromLastGame(lastGame *Game) *Game {
	winners := make([]*Player, 0, lastGame.MaxWinners)

	newGame := &Game{
		TotalCells:  lastGame.TotalCells,
		MaxPlayers:  lastGame.MaxPlayers,
		MaxWinners:  lastGame.MaxWinners,
		DiceCount:   lastGame.DiceCount,
		GameBoard:   NewBoard(lastGame.TotalCells),
		Status:      Created,
		CurrentTurn: 0,
		DiceManager: lastGame.DiceManager,
		Creator:     lastGame.Creator,
		Winners:     winners,
	}

	// add all players from last game
	for _, p := range lastGame.Players {
		newGame.AddUser(&User{
			Id:      p.Id,
			Name:    p.Name,
			Picture: p.Picture,
		})
	}

	newGame.Messages = []utils.ChatMessage{
		getSystemMessage("Game created from last game!"),
	}

	return newGame
}

func getSystemMessage(message string) utils.ChatMessage {
	return utils.ChatMessage{
		PlayerName: "System",
		PlayerId:   "system",
		Message:    message,
	}
}
