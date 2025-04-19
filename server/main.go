package main

import (
	"log"

	"github.com/sourcecode081017/im-chat-golang-react/websocket"
)

func main() {
	// Start WebSocket server
	hub := websocket.NewHub()
	// start hub server in a goroutine
	log.Println("Starting hub server...")
	go hub.Run()
	// create a websocker object
	ws := websocket.NewWs(hub, "client1")
	// start the websocket server
	ws.StartWebSocketServer()

}
