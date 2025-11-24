package handle

import (
	"net/http"

	"gibhub.com/l-jessie/test-im/internal/logic"
	"gibhub.com/l-jessie/test-im/internal/model/dto"
	"gibhub.com/l-jessie/test-im/internal/model/entity"
	"gibhub.com/l-jessie/test-im/internal/model/types"
	"gibhub.com/l-jessie/test-im/internal/utils"

	"github.com/gin-gonic/gin"
)

type RoomHandle struct {
	hub         *types.Hub
	chatService *logic.ChatService
}

func NewRoomHandle(hub *types.Hub, chatService *logic.ChatService) *RoomHandle {
	return &RoomHandle{hub: hub, chatService: chatService}
}

func (h *RoomHandle) GetRoomsHandle(c *gin.Context) {
	rooms := make([]*dto.RoomsResponse, 0, len(h.hub.Rooms))
	for _, room := range h.hub.Rooms {
		rooms = append(rooms, &dto.RoomsResponse{
			ID:          room.ID,
			Name:        room.Name,
			HasPassword: room.Password != "",
			UserID:      room.UserID,
			UserName:    room.UserName,
			Count:       len(room.Clients),
			CreateTime:  entity.BizTimeFull(room.CreateTime),
		})
	}

	c.JSON(http.StatusOK,
		gin.H{"code": 1, "msg": "success", "data": rooms},
	)
}

func (h *RoomHandle) CreateRoomHandle(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	var oneClient types.Client
	userCli := h.hub.Users[req.UserID]
	if len(userCli) >= 1 {
		for cli := range userCli {
			oneClient = *cli
			break
		}
	} else {
		oneClient = types.Client{}
	}

	roomID := utils.GenerateUUID()
	room := types.NewRoom(roomID, req.Name, req.Password, req.UserID, oneClient.UserName)

	h.hub.CreateRoom <- &types.CreateRoomEvent{
		UserID: req.UserID,
		RoomID: roomID,
		Room:   room,
	}

	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &dto.RoomsResponse{
		ID:          room.ID,
		Name:        room.Name,
		HasPassword: room.Password != "",
		UserID:      room.UserID,
		UserName:    room.UserName,
		Count:       len(room.Clients),
		CreateTime:  entity.BizTimeFull(room.CreateTime),
	}})
}

func (h *RoomHandle) JoinRoomHandle(c *gin.Context) {
	var req dto.JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	if rooms, ok := h.hub.Rooms[req.RoomID]; ok {
		if rooms.Password != "" && req.Password != rooms.Password {
			c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "密码错误"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "房间不存在"})
		return
	}

	h.hub.JoinRoom <- &types.JoinRoomEvent{
		UserID:   req.UserID,
		DeviceID: req.DeviceID,
		RoomID:   req.RoomID,
	}

	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success"})
}

func (h *RoomHandle) GetRoomDetailHandle(c *gin.Context) {
	roomID := c.Param("roomId")

	room, ok := h.hub.Rooms[roomID]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "msg": "房间不存在"})
		return
	}

	users := make([]*dto.UserVO, 0, len(room.Clients))
	i := 0
	for client := range room.Clients {
		i++
		users = append(users, &dto.UserVO{
			ID:   client.UserId,
			Name: client.UserName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &dto.RoomsDetailResponse{
		ID:         room.ID,
		Name:       room.Name,
		UserID:     room.UserID,
		UserName:   "room-admin",
		Count:      len(room.Clients),
		CreateTime: entity.BizTimeFull(room.CreateTime),
		Users:      users,
	}})
}
