package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sourcecode081017/im-chat-golang-react/db/postgres"
	"github.com/sourcecode081017/im-chat-golang-react/websocket"
)

func main() {

	// load configs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	// create a context to be used throughout the application
	ctx := context.Background()
	// create a new postgres database connection
	pgDb, err := postgres.NewPgDb(ctx)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
		os.Exit(1)
	}
	fmt.Println(pgDb)
	// run the database migrations
	err = pgDb.RunMigrations(ctx)
	if err != nil {
		log.Fatal("Error running migrations: ", err)
		os.Exit(1)
	}
	// Start WebSocket server
	hub := websocket.NewHub()
	// start hub server in a goroutine
	log.Println("Starting hub server...")
	go hub.Run()
	// create a websocker object
	ws := websocket.NewWs(hub, pgDb)
	// start the websocket server
	ws.StartWebSocketServer()

}
