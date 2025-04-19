package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Ws struct {
	wsHub    *Hub
	clientId string
}

func NewWs(hub *Hub, clientId string) *Ws {
	return &Ws{
		wsHub:    hub,
		clientId: clientId,
	}
}

func (ws *Ws) StartWebSocketServer() {
	// Initialize a http server using net/http
	router := gin.Default()
	// upgrade the http connection request to a websocket connection
	router.GET("/connect", ws.serveWebsocket)
	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic("Failed to start WebSocket server: " + err.Error())
	}
	// Log the server start
	log.Println("WebSocket server started on :8080")
}

// serveWebsocket handles the WebSocket connection
func (ws *Ws) serveWebsocket(c *gin.Context) {
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
	client := NewClient(ws.clientId, conn, ws.wsHub)
	client.hub.register <- client
	// Start the read and write pumps
	go client.ReadPump()

	go client.WritePump()

	defer conn.Close()
}
