package handle

import (
	"github.com/l-jessie/test-im/internal/model/dto"
	"github.com/l-jessie/test-im/internal/model/types"

	"github.com/gin-gonic/gin"
)

type UsersHandle struct {
	hub *types.Hub
}

func NewUsersHandle(hub *types.Hub) *UsersHandle {
	return &UsersHandle{
		hub: hub,
	}
}

func (h *UsersHandle) GetUsersHandle(context *gin.Context) {
	users := h.hub.Users

	userInfos := make([]*dto.UserVO, 0, len(users))
	clients := h.hub.Clients
	for cli := range clients {
		userInfos = append(userInfos, &dto.UserVO{
			ID:   cli.UserId,
			Name: cli.UserName,
		})
	}

	context.JSON(200, gin.H{
		"code": 1,
		"msg":  "success",
		"data": userInfos,
	})
}
