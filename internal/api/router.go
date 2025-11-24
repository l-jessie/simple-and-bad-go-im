package api

import (
	"github.com/l-jessie/test-im/internal/handle"
	"github.com/l-jessie/test-im/internal/logic"
	"github.com/l-jessie/test-im/internal/model/types"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// DI
	hub := types.NewHub()
	go hub.Run()
	chatService := logic.NewChatService(hub)
	wsHandle := handle.NewWsHandle(hub, chatService)
	roomHandle := handle.NewRoomHandle(hub, chatService)
	usersHandle := handle.NewUsersHandle(hub)

	// 路由
	router := gin.Default()
	v1Group := router.Group("/v1/api")
	v1Group.GET("/ping", handle.PingPongHandle)
	v1Group.POST("/login", handle.LoginHandle)
	v1Group.GET("/ws", wsHandle.WsHandleFunc)

	roomGroup := v1Group.Group("/rooms")
	{
		roomGroup.GET("", roomHandle.GetRoomsHandle)
		roomGroup.POST("", roomHandle.CreateRoomHandle)
		roomGroup.GET("/:roomId", roomHandle.GetRoomDetailHandle)
		roomGroup.POST("/:roomId/join", roomHandle.JoinRoomHandle)
	}

	usersGroup := v1Group.Group("users")
	{
		usersGroup.GET("", usersHandle.GetUsersHandle)
	}

	return router
}
