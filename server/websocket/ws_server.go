package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func StartWebSocketServer() {
	// Initialize a http server using net/http
	router := gin.Default()
	// upgrade the http connection request to a websocket connection
	router.GET("/ws", serveWebsocket)
	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic("Failed to start WebSocket server: " + err.Error())
	}
	// Log the server start
	log.Println("WebSocket server started on :8080")
}

// serveWebsocket handles the WebSocket connection
func serveWebsocket(c *gin.Context) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return false
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()
}
