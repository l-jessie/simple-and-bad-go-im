package handle

import (
	"net/http"
	"strings"

	"gibhub.com/l-jessie/test-im/internal/logic"
	types2 "gibhub.com/l-jessie/test-im/internal/model/types"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandle struct {
	hub         *types2.Hub
	chatService *logic.ChatService
	upgrader    websocket.Upgrader
}

func matchOrigin(pattern, origin string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(origin, strings.TrimSuffix(pattern, "*"))
	}
	return pattern == origin
}

func NewWsHandle(hub *types2.Hub, chatService *logic.ChatService) *WsHandle {
	return &WsHandle{
		hub:         hub,
		chatService: chatService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")

				allowed := []string{
					"http://localhost",
					"http://localhost:5173",
					"https://*.your-company.com",
				}

				for _, a := range allowed {
					if matchOrigin(a, origin) {
						return true
					}
				}

				return false
			},
		},
	}
}

func (w *WsHandle) WsHandleFunc(c *gin.Context) {
	userID := c.Query("token")
	deviceID := c.Query("deviceId")
	userName := c.Query("username")

	if userID == "" || deviceID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 0,
			"msg":  "参数错误",
		})
		return
	}

	connect, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  "升级失败",
		})
		return
	}

	client := types2.NewClient(w.hub, connect, userID, userName, deviceID)
	w.hub.Register <- &types2.RegisterEvent{
		Client: client,
		UserId: userID,
	}

	// Start read and write pumps
	go client.WritePump()
	go client.ReadPump(w.chatService.HandleMessage)
}
