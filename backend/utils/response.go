package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type WSErrorRes struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SendHTTPErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	resp := &Response{
		Status:  statusCode,
		Message: message,
		Data:    nil,
	}

	json.NewEncoder(w).Encode(resp)
}

func SendHTTPSuccessResponse(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)

	resp := &Response{
		Status:  statusCode,
		Message: "success",
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

func ErrorResponse(statusCode int, message string) WSErrorRes {
	return WSErrorRes{
		Status:  statusCode,
		Message: message,
	}
}

type ChatMessage struct {
	PlayerName string `json:"player_name"`
	PlayerId   string `json:"player_id"`
	Message    string `json:"message"`
}

type NextGameMessage struct {
	GameId string `json:"game_id"`
}