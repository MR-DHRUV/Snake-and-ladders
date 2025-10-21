package rest

import (
	"encoding/json"
	"net/http"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/MR-DHRUV/snake_and_ladders/constants"
	"github.com/MR-DHRUV/snake_and_ladders/services"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
)

type AuthRequest struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}

func GoogleAuthController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tokenReq AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if tokenReq.CodeVerifier == "" || tokenReq.Code == "" {
		utils.SendHTTPErrorResponse(w, http.StatusBadRequest, "code and code_verifier are required")
		return
	}

	user, err := services.GoogleAuthService(tokenReq.Code, tokenReq.CodeVerifier)
	if err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusInternalServerError, "authentication failed: "+err.Error())
		return
	}

	token, err := utils.GenerateJWT(user.Id)
	if err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusInternalServerError, "failed to generate token: "+err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constants.AuthToken,
		Value:    token,
		Path:     "/",                     // cookie is valid for entire site
		HttpOnly: true,                    // can't be accessed by JS
		Secure:   config.AuthCookieSecure, // only over HTTPS
		SameSite: http.SameSiteLaxMode,    // prevents CSRF in most cases
		MaxAge:   config.AuthCookieExpiry, // 1 week
	})

	response := map[string]interface{}{
		"success": true,
	}

	utils.SendHTTPSuccessResponse(w, http.StatusOK, response)
}

func GetUserController(w http.ResponseWriter, r *http.Request) {
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

	user, err := services.GetUserById(userId)
	if err != nil {
		utils.SendHTTPErrorResponse(w, http.StatusInternalServerError, "failed to get user: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"user": user,
	}

	utils.SendHTTPSuccessResponse(w, http.StatusOK, response)
}

func HealthCheckController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": "ok",
	}
	utils.SendHTTPSuccessResponse(w, http.StatusOK, response)
}