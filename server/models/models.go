package models

type Message struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	ChannelId   string `json:"channelId"`
	RecipientId string `json:"recipientId"`
}
