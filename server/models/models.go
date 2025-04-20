package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	ChannelId   string `json:"channelId"`
	RecipientId string `json:"recipientId"`
}

type Channel struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Subscribers []User    `json:"subscribers"`
}
type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}
