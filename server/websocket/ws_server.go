package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sourcecode081017/im-chat-golang-react/db/postgres"
	"github.com/sourcecode081017/im-chat-golang-react/models"
)

type Ws struct {
	wsHub *Hub
	pgDb  *postgres.PgDb
}

func NewWs(hub *Hub, pgDb *postgres.PgDb) *Ws {
	return &Ws{
		wsHub: hub,
		pgDb:  pgDb,
	}
}

func (ws *Ws) StartWebSocketServer() {
	// Initialize a http server using net/http
	router := gin.Default()
	// upgrade the http connection request to a websocket connection
	router.GET("/connect/:clientId", ws.serveWebsocket)
	// route to create a new user channel
	router.POST("/user/:userId/channel", ws.createUserChannel)
	// route to create a new user
	router.POST("/user", ws.createUser)
	// route to fetch all users
	router.GET("/users", ws.fetchAllUsers)
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

func (ws *Ws) createUser(c *gin.Context) {
	// Implement the logic to create a user
	// For example, you can use the userId to create a new user in your database
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid input",
		})
		return
	}
	// Create the user in the database
	if err := ws.pgDb.CreateUser(c, &user); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "User created successfully",
	})
}

func (ws *Ws) fetchAllUsers(c *gin.Context) {
	// Implement the logic to fetch all users
	// For example, you can use the userId to fetch all users from your database
	users, err := ws.pgDb.GetUsers(c)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}
	c.JSON(200, gin.H{
		"users": users,
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
	// Create a new client
	client := NewClient(clientId, conn, ws.wsHub)
	client.hub.register <- client
	// Start the read and write pumps
	go client.ReadPump()

	go client.WritePump()

	//defer conn.Close()
}
