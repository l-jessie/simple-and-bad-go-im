package types

import (
	"errors"
	"log"
	"time"

	"github.com/l-jessie/test-im/internal/global"

	"github.com/gorilla/websocket"
)

// Client 是 websocket 连接和 hub 之间的中间人。
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn // 从 WsConnect 重命名而来，为了简洁和一致性
	send     chan []byte
	UserId   string
	UserName string
	DeviceId string
}

func NewClient(hub *Hub, conn *websocket.Conn, userId, userName, deviceId string) *Client {
	return &Client{
		Hub:      hub,
		Conn:     conn,
		send:     make(chan []byte, 256),
		UserId:   userId,
		UserName: userName,
		DeviceId: deviceId,
	}
}

// readPump 将消息从 websocket 连接泵送到 hub。
//
// 应用程序在每个连接的 goroutine 中运行 readPump。应用程序
// 通过从此 goroutine 执行所有读取来确保连接上最多只有一个读取器。
func (c *Client) ReadPump(messageFunc func(client *Client, message []byte)) {
	defer func() {
		// 退出时，注销客户端并关闭连接
		c.Hub.Unregister <- &UnRegisterEvent{Client: c, UserId: c.UserId}
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(global.MaxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(global.PongWait))
	c.Conn.SetPongHandler(func(string) error { _ = c.Conn.SetReadDeadline(time.Now().Add(global.PongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: UserID: %s, %v", c.UserId, err)
			}
			break // 读取错误时退出循环
		}
		// 将原始消息传递给处理函数
		messageFunc(c, message)
	}
}

// writePump 将消息从 hub 泵送到 websocket 连接。
//
// 为每个连接启动一个运行 writePump 的 goroutine。应用程序
// 通过从此 goroutine 执行所有写入来确保连接上最多只有一个写入器。
func (c *Client) WritePump() {
	ticker := time.NewTicker(global.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(global.WriteWait))
			if !ok {
				// hub 关闭了 channel。
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("write error: UserID: %s, %v", c.UserId, err)
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(global.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ping error: UserID: %s, %v", c.UserId, err)
				return
			}
		}
	}
}

// SendMessage 是向客户端发送消息的线程安全方式。
// 如果发送 channel 已关闭，它会从 panic 中恢复。
func (c *Client) SendMessage(message []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic in SendMessage for UserID %s: %v", c.UserId, r)
			err = errors.New("send to closed channel") // 设置一个有意义的错误
		}
	}()
	// 非阻塞发送以防止阻塞 hub 的广播。
	select {
	case c.send <- message:
	default:
		log.Printf("send channel full for UserID: %s. Message dropped.", c.UserId)
		err = errors.New("send channel full") // 设置一个有意义的错误
	}
	return err // 返回错误
}
