package types

import (
	"encoding/json"
)

type MessageEventType int

const (
	ReloadUsers MessageEventType = iota
	ReloadRoomsDetail
	ReloadRooms
)

type MessageEvent struct {
	Type MessageEventType `json:"type"`
	Data json.RawMessage  `json:"data"`
}

func NewMessageEventPayload(t MessageEventType, data []byte) *MessageEvent {
	return &MessageEvent{
		Type: t,
		Data: data,
	}
}
