package routes

import (
	"net/http"

	controller "github.com/MR-DHRUV/snake_and_ladders/controller/rest"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func RegisterHTTPRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/auth/google", controller.GoogleAuthController).Methods("POST")
	router.HandleFunc("/auth/user", controller.GetUserController).Methods("GET")
	router.HandleFunc("/game", controller.CreateGameController).Methods("POST")
	router.HandleFunc("/past-games", controller.GetGamesController).Methods("GET")
	router.HandleFunc("/health", controller.HealthCheckController).Methods("GET")

	return router
}

func StartHTTPServer() {
	router := RegisterHTTPRoutes()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5000", "https://snake-and-ladders-ten.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)
	utils.GetLogger().Info("HTTP server is running on port 8081")

	if err := http.ListenAndServe(":8081", handler); err != nil {
		utils.GetLogger().Error("Failed to start HTTP server: %v", err)
	}
}
