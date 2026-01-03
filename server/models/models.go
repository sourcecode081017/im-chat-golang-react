package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message represents the wire format for WebSocket communication
type Message struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	ChannelId   string `json:"channelId"`
	RecipientId string `json:"recipientId"`
	SenderId    string `json:"senderId,omitempty"` // Added for sender identification
}

// ChatMessage represents a persisted message in the database
type ChatMessage struct {
	gorm.Model
	SenderID    uuid.UUID `json:"senderId" gorm:"type:uuid;not null;index"`
	RecipientID uuid.UUID `json:"recipientId" gorm:"type:uuid;not null;index"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	MessageType string    `json:"messageType" gorm:"type:varchar(50);default:'direct'"`
}

type Channel struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Subscribers []User    `json:"subscribers"`
}
type User struct {
	gorm.Model
	Username  string    `json:"username" gorm:"unique; not null"`
	UserUUID  uuid.UUID `json:"userId" gorm:"type:uuid; unique; not null"`
	Email     *string   `json:"email"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
}
