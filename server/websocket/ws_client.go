package websocket

import (
	"log"

	"github.com/gorilla/websocket"
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
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for message := range c.send {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c *Client) SendMessage(message []byte) {
	c.send <- message
}

func (c *Client) Close() {
	c.hub.unregister <- c
	close(c.send)
	c.conn.Close()
}
