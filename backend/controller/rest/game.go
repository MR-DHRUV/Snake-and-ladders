package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MR-DHRUV/snake_and_ladders/constants"
	"github.com/MR-DHRUV/snake_and_ladders/services"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
)

type CreateGameRequest struct {
	MaxPlayers int    `json:"max_players"`
	MaxWinners int    `json:"max_winners"`
	DiceCount  int    `json:"dice_count"`
}

func CreateGameController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie(constants.AuthToken)
	if err != nil {
		http.Error(w, "Authorization cookie missing", http.StatusUnauthorized)
		return
	}

	_, userId, err := utils.VerifyJWT(cookie.Value)
	if err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusUnauthorized, "invalid token: "+err.Error())
		return
	}

	var createGameReq CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&createGameReq); err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	game, err := services.CreateNewGame(
		userId,
		createGameReq.MaxPlayers,
		createGameReq.MaxWinners,
		createGameReq.DiceCount,
	)

	if err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusInternalServerError, "failed to create game: "+err.Error())
		return
	}

	utils.SendHTTPSuccessResponse(w, http.StatusOK, game.Id)
}

func GetGamesController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie(constants.AuthToken)
	if err != nil {
		http.Error(w, "Authorization cookie missing", http.StatusUnauthorized)
		return
	}

	_, userId, err := utils.VerifyJWT(cookie.Value)
	if err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusUnauthorized, "invalid token: "+err.Error())
		return
	}

	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err == nil && p > 0 {
			page = p
		}
	}

	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err == nil && l > 0 {
			limit = l
		}
	}

	games, err := services.GetGames(userId, page, limit)
	if err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusInternalServerError, "failed to fetch games: "+err.Error())
		return
	}

	utils.SendHTTPSuccessResponse(w, http.StatusOK, games)
}
