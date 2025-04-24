package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	ChannelId   string `json:"channelId"`
	RecipientId string `json:"recipientId"`
}

type User struct {
	gorm.Model
	Username  string    `json:"username" gorm:"unique; not null"`
	UserUUID  uuid.UUID `json:"userId" gorm:"type:uuid; unique; not null"`
	Email     *string   `json:"email"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
}

type Channel struct {
	gorm.Model
	ChannelUUID uuid.UUID `json:"channelId" gorm:"type:uuid; unique; not null"`
	ChannelName string    `json:"channelName" gorm:"not null"`
	CreatedBy   string    `json:"createdBy" gorm:"not null"`
	Users       []*User   `gorm:"many2many:user_channels"`
}
