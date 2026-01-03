package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sourcecode081017/im-chat-golang-react/db/postgres"
	"github.com/sourcecode081017/im-chat-golang-react/models"
)

type Client struct {
	Id   string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	pgDb *postgres.PgDb
}

func NewClient(id string, conn *websocket.Conn, hub *Hub, pgDb *postgres.PgDb) *Client {
	return &Client{
		Id:   id,
		conn: conn,
		send: make(chan []byte),
		hub:  hub,
		pgDb: pgDb,
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

		// Enrich message with sender information
		msg.SenderId = c.Id

		// Parse UUIDs for database storage
		senderUUID, err := uuid.Parse(c.Id)
		if err != nil {
			log.Printf("invalid sender UUID: %v", err)
			continue
		}

		recipientUUID, err := uuid.Parse(msg.RecipientId)
		if err != nil {
			log.Printf("invalid recipient UUID: %v", err)
			continue
		}

		// Save to database
		chatMessage := &models.ChatMessage{
			SenderID:    senderUUID,
			RecipientID: recipientUUID,
			Content:     msg.Content,
			MessageType: msg.Type,
		}

		if err := c.pgDb.SaveMessage(context.Background(), chatMessage); err != nil {
			log.Printf("failed to save message: %v", err)
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
