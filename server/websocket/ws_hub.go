package websocket

import (
	"encoding/json"
	"log"

	"github.com/sourcecode081017/im-chat-golang-react/models"
)

type Hub struct {
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Id] = client
			log.Printf("Client %s connected", client.Id)
		case client := <-h.unregister:
			if _, ok := h.clients[client.Id]; ok {
				delete(h.clients, client.Id)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Printf("Broadcasting message: %s", message)
			var msg *models.Message
			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			// for one on one chat
			for id, client := range h.clients {
				if msg.RecipientId != "" && msg.RecipientId != id {
					continue
				}
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, id)
				}
			}
		}
	}
}
