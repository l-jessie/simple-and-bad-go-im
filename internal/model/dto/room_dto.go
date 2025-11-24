package dto

import (
	"test-im/internal/model/entity"
)

type CreateRoomRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	UserID   string `json:"userId"` // 房间拥有者 用户ID
}

type JoinRoomRequest struct {
	RoomID   string `json:"roomId"`
	Password string `json:"password"`
	UserID   string `json:"userId"`
	DeviceID string `json:"deviceId"`
}

type RoomsResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	HasPassword bool               `json:"hasPassword"` // 是否有密码
	UserID      string             `json:"userId"`      // 房间拥有者 用户ID
	UserName    string             `json:"userName"`    // 房间拥有者 用户名称
	Count       int                `json:"count"`       // 房间内用户数量
	CreateTime  entity.BizTimeFull `json:"createTime"`  // 创建时间
}

type RoomsDetailResponse struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	UserID     string             `json:"userId"`     // 房间拥有者 用户ID
	UserName   string             `json:"userName"`   // 房间拥有者 用户名称
	Count      int                `json:"count"`      // 房间内用户数量
	CreateTime entity.BizTimeFull `json:"createTime"` // 创建时间
	Users      []*UserVO          `json:"users"`
}
