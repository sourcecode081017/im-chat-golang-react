package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Ws struct {
	wsHub *Hub
}

func NewWs(hub *Hub) *Ws {
	return &Ws{
		wsHub: hub,
	}
}

func (ws *Ws) StartWebSocketServer() {
	// Initialize a http server using net/http
	router := gin.Default()
	// upgrade the http connection request to a websocket connection
	router.GET("/connect/:clientId", ws.serveWebsocket)
	// route to create a new user channel
	router.POST("/user/:userId/channel", ws.createUserChannel)
	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic("Failed to start WebSocket server: " + err.Error())
	}
	// Log the server start
	log.Println("WebSocket server started on :8080")
}

func (ws *Ws) createUserChannel(c *gin.Context) {
	userId := c.Param("userId")
	// Implement the logic to create a user channel
	// For example, you can use the userId to create a new channel in your database
	c.JSON(200, gin.H{
		"message": "User channel created successfully",
		"userId":  userId,
	})
}

// serveWebsocket handles the WebSocket connection
func (ws *Ws) serveWebsocket(c *gin.Context) {
	clientId := c.Param("clientId")
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// create a

	// Create a new client
	client := NewClient(clientId, conn, ws.wsHub)
	client.hub.register <- client
	// Start the read and write pumps
	go client.ReadPump()

	go client.WritePump()

	//defer conn.Close()
}
