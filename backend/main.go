package main

import (
    "os"
    "os/signal"
    "sync"
    "syscall"

	"github.com/MR-DHRUV/snake_and_ladders/db"
	"github.com/MR-DHRUV/snake_and_ladders/routes"
	"github.com/MR-DHRUV/snake_and_ladders/utils"
	"github.com/MR-DHRUV/snake_and_ladders/transport"
	controller "github.com/MR-DHRUV/snake_and_ladders/controller/ws"
)

func main() {
	// Init Db connections
	db.InitMongo()
	db.InitRedis()

	utils.GetLogger().Info("Connected to MongoDB")

	var wg sync.WaitGroup

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.GetLogger().Info("Starting HTTP server...")
		routes.StartHTTPServer()
	}()

	// Start WebSocket server
	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.GetLogger().Info("Starting WebSocket server...")
		transport.GetConnectionManager(); // init conncection manager
		routes.StartWebSocketServer();
	}()

	// Start Redis Subscriber
	wg.Add(1)
	go func() {
		defer wg.Done()
		utils.GetLogger().Info("Starting Redis Subscriber...")
		controller.StartGameEventSubscriber()
	}()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1) // create a channel to listen for OS signals of size 1
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM) // listen for interrupt and terminate signals

	go func() {
		sig := <-sigCh // block until a signal is received
		// gracefully shutdown servers
		utils.GetLogger().Info("Received signal: %v, shutting down...", sig) 
		os.Exit(0)
	}()

	wg.Wait()
	utils.GetLogger().Error("All servers have stopped. Exiting...")
}
