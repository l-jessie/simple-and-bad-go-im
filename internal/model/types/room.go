package types

import (
	"time"
)

type Room struct {
	ID         string
	Name       string
	Password   string
	UserID     string // 房间拥有者 用户ID
	UserName   string
	Clients    map[*Client]bool
	CreateTime time.Time
}

func NewRoom(id, name, password, ownerUserID, ownerUserName string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Password:   password,
		Clients:    make(map[*Client]bool),
		UserID:     ownerUserID,
		UserName:   ownerUserName,
		CreateTime: time.Now(),
	}
}
