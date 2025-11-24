package logic

import (
	"encoding/json"
	"log"

	types2 "gibhub.com/l-jessie/test-im/internal/model/types"
)

type ChatService struct {
	hub *types2.Hub
}

func NewChatService(hub *types2.Hub) *ChatService {
	return &ChatService{
		hub: hub,
	}
}

func (c *ChatService) HandleMessage(client *types2.Client, messageByte []byte) {
	var message *types2.Message
	err := json.Unmarshal(messageByte, &message)
	if err != nil {
		log.Printf("message unmarshal error: %v", err)
		return
	}

	message.From = client.UserId
	c.hub.Broadcast <- message
}
