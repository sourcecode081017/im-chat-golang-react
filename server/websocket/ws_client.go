package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/sourcecode081017/im-chat-golang-react/models"
)

type Client struct {
	Id   string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Id:   id,
		conn: conn,
		send: make(chan []byte),
		hub:  hub,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var msg *models.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		c.hub.broadcast <- msgJSON
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
