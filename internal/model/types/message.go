package types

import (
	"time"
)

type MessageType int

const (
	MessageTypeJoin MessageType = iota
	MessageTypeLeave
	MessageTypeRoom
	MessageTypeUser
	MessageTypeGlobal
	MessageTypeSystem
	MessageTypeJoinRoom
	MessageTypeLeaveRoom
)

type Message struct {
	Type         MessageType   `json:"type"`
	Payload      *Payload      `json:"payload"`
	From         string        `json:"from"`
	To           string        `json:"to"`
	MessageEvent *MessageEvent `json:"messageEvent"`
	Timestamp    int64         `json:"time,omitempty"`
}

func NewMessage(t MessageType, payload *Payload, from string, to string) *Message {
	return &Message{
		Type:      t,
		Payload:   payload,
		From:      from,
		To:        to,
		Timestamp: time.Now().Unix(),
	}
}

func NewMessageEvent(t MessageType, messageEvent *MessageEvent) *Message {
	return &Message{
		Type:         t,
		MessageEvent: messageEvent,
		Timestamp:    time.Now().Unix(),
	}
}
