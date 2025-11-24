package types

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
)

var (
	ClientNotFindError = errors.New("client not find")
)

type Hub struct {
	// mu 保护下面的 map。
	mu sync.RWMutex

	Clients   map[*Client]bool              // 所有链接的 Set
	Users     map[string]map[*Client]bool   // 用户所有链接的 Set  userId -> 链接
	UserInfos map[string]map[*UserInfo]bool // 用户信息合集 userId -> 用户信息

	Rooms     map[string]*Room           // 房间ID, 房间
	UserRooms map[string]map[string]bool // 用户ID, 房间IDs

	Broadcast  chan *Message         // 广播
	Register   chan *RegisterEvent   // 注册链接
	Unregister chan *UnRegisterEvent // 注销链接

	CreateRoom chan *CreateRoomEvent // 创建房间
	JoinRoom   chan *JoinRoomEvent   // 加入房间
	UnjoinRoom chan *UnJoinRoomEvent // 退出房间
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]map[*Client]bool),
		UserInfos:  make(map[string]map[*UserInfo]bool),
		Rooms:      make(map[string]*Room),
		UserRooms:  make(map[string]map[string]bool),
		Broadcast:  make(chan *Message),
		CreateRoom: make(chan *CreateRoomEvent),
		Register:   make(chan *RegisterEvent),
		Unregister: make(chan *UnRegisterEvent),
		JoinRoom:   make(chan *JoinRoomEvent),
		UnjoinRoom: make(chan *UnJoinRoomEvent),
	}
}

func (h *Hub) Run() {
	for {
		select {
		// 注册链接
		case event := <-h.Register:
			registerClient(h, event)

		// 注销链接
		case event := <-h.Unregister:
			unRegisterClient(h, event)

		// 广播
		case message := <-h.Broadcast:
			broadcast(h, message)

		// 创建房间
		case createRoom := <-h.CreateRoom:
			createRoomFunc(h, createRoom)

		// 加入房间(不存在创建房间)
		case event := <-h.JoinRoom:
			joinRoom(h, event)

		// 退出房间
		case event := <-h.UnjoinRoom:
			unjoinRoom(h, event)
		}
	}
}

func registerClient(h *Hub, event *RegisterEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("register client: %s", event.UserId)

	// 添加到全局客户端集合
	h.Clients[event.Client] = true

	// 添加到用户特定的客户端集合
	if _, ok := h.Users[event.UserId]; !ok {
		h.Users[event.UserId] = make(map[*Client]bool)
	}
	h.Users[event.UserId][event.Client] = true

	// 通知用户上线了刷新 /users接口
	go func() {
		marshal, _ := json.Marshal(map[string]string{
			"userId":  event.UserId,
			"message": "online",
		})

		h.Broadcast <- NewMessageEvent(
			MessageTypeGlobal,
			NewMessageEventPayload(ReloadUsers, marshal),
		)
	}()
}

func unRegisterClient(h *Hub, event *UnRegisterEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 检查客户端是否存在以防止重复注销
	if _, ok := h.Clients[event.Client]; ok {
		// 关闭发送通道以停止向其发送消息
		close(event.Client.send)

		// 从全局客户端集合中删除
		delete(h.Clients, event.Client)

		// 从用户特定的客户端集合中删除
		if clients, ok := h.Users[event.UserId]; ok {
			delete(clients, event.Client)
			if len(clients) == 0 {
				delete(h.Users, event.UserId)
			}
		}

		// 将客户端从其所在的任何房间中删除
		for roomId, room := range h.Rooms {
			if _, ok := room.Clients[event.Client]; ok {
				delete(room.Clients, event.Client)
				// 如果客户端离开后房间为空，则删除房间
				if len(room.Clients) == 0 {
					delete(h.Rooms, roomId)
					// 同时从 UserRooms map 中删除以确保彻底
					if userRoomSet, ok := h.UserRooms[room.UserID]; ok {
						delete(userRoomSet, roomId)
						if len(userRoomSet) == 0 {
							delete(h.UserRooms, room.UserID)
						}
					}
				}
			}
		}

		go func() {
			marshal, _ := json.Marshal(map[string]string{
				"userId":  event.UserId,
				"message": "offline",
			})
			h.Broadcast <- NewMessageEvent(
				MessageTypeGlobal,
				NewMessageEventPayload(ReloadUsers, marshal),
			)
		}()

		log.Printf("unregistered client: %s", event.UserId)
	}
}

func broadcast(h *Hub, msg *Message) {
	messageMarshal, err := json.Marshal(msg)
	if err != nil {
		log.Printf("broadcast json marshal error: %v", err)
		return
	}

	h.mu.RLock()         // 获取读锁
	defer h.mu.RUnlock() // 函数退出时释放读锁

	if msg == nil {
		log.Printf("broadcast message is nil")
		return
	}

	var targetClients []*Client // 收集要发送到的客户端
	if msg.Type == MessageTypeGlobal {
		for client := range h.Clients {
			targetClients = append(targetClients, client)
		}
	} else if msg.Type == MessageTypeRoom {
		if room, ok := h.Rooms[msg.To]; ok {
			for client := range room.Clients {
				targetClients = append(targetClients, client)
			}
		}
	} else if msg.Type == MessageTypeUser {
		if clients, ok := h.Users[msg.To]; ok {
			for client := range clients {
				targetClients = append(targetClients, client)
			}
		}
	}

	// 现在在锁范围之外使用收集到的 targetClients 发送消息
	// 注意：SendMessage 是非阻塞的，并且会处理自己的 panic，因此在锁外调用是安全的。
	for _, client := range targetClients {
		err := client.SendMessage(messageMarshal)
		if err != nil {
			log.Printf("broadcast SendMessage error: %v", err)
			return
		}
	}
}

func createRoomFunc(h *Hub, event *CreateRoomEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("create room: %s", event.Room.ID)

	h.Rooms[event.Room.ID] = event.Room

	// 添加到用户的房间映射中
	if _, ok := h.UserRooms[event.Room.UserID]; !ok {
		h.UserRooms[event.Room.UserID] = make(map[string]bool)
	}
	h.UserRooms[event.Room.UserID][event.Room.ID] = true

	go func() {
		h.Broadcast <- NewMessageEvent(
			MessageTypeGlobal,
			NewMessageEventPayload(ReloadRooms, nil),
		)
	}()
}

func joinRoom(h *Hub, event *JoinRoomEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("join room: %s", event.RoomID)

	room, ok := h.Rooms[event.RoomID]
	if !ok {
		log.Printf("room not exist: %s", event.RoomID)
		return
	}

	// 调用 NoLock 版本，因为我们已经持有了写锁
	currentClient, err := h.findClientNoLock(event.UserID, event.DeviceID)
	if errors.Is(err, ClientNotFindError) {
		log.Printf("client not exist: %s", event.UserID)
		return
	}

	// 将客户端添加到房间的客户端集合中
	// 确保 room.Clients 已初始化。如果 NewRoom 被正确调用，它应该是初始化的。
	// 但是，检查更安全。Room 结构体有 Clients map。
	if room.Clients == nil {
		room.Clients = make(map[*Client]bool)
	}
	room.Clients[currentClient] = true

	// 将房间添加到用户的房间映射中
	if _, ok := h.UserRooms[event.UserID]; !ok {
		h.UserRooms[event.UserID] = make(map[string]bool)
	}
	h.UserRooms[event.UserID][event.RoomID] = true

	// Broadcast ReloadRoomsDetail to all clients
	roomIDBytes, err := json.Marshal(event.RoomID)
	if err != nil {
		log.Printf("Error marshalling room ID for ReloadRoomsDetail: %v", err)
	} else {
		go func() {
			h.Broadcast <- NewMessageEvent(
				MessageTypeGlobal, // This should be MessageTypeGlobal to be caught by the frontend as a global event
				NewMessageEventPayload(ReloadRoomsDetail, json.RawMessage(roomIDBytes)),
			)
		}()

		go func() {
			h.Broadcast <- NewMessageEvent(
				MessageTypeGlobal,
				NewMessageEventPayload(ReloadRooms, nil),
			)
		}()
	}
}

func unjoinRoom(h *Hub, event *UnJoinRoomEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Printf("unjoin room: %s for user: %s, device: %s", event.RoomID, event.UserID, event.DeviceID)

	room, ok := h.Rooms[event.RoomID]
	if !ok {
		log.Printf("room not exist: %s", event.RoomID)
		return
	}

	currentClient, err := h.findClientNoLock(event.UserID, event.DeviceID)
	if errors.Is(err, ClientNotFindError) {
		log.Printf("client not exist: %s", event.UserID)
		return
	}

	// Remove client from room's clients
	if _, ok := room.Clients[currentClient]; ok {
		delete(room.Clients, currentClient)
	}

	// Remove room from user's rooms
	if userRoomSet, ok := h.UserRooms[event.UserID]; ok {
		delete(userRoomSet, event.RoomID)
		if len(userRoomSet) == 0 {
			delete(h.UserRooms, event.UserID)
		}
	}

	// If the room becomes empty, remove the room itself
	if len(room.Clients) == 0 {
		delete(h.Rooms, event.RoomID)
		log.Printf("room %s is now empty and removed", event.RoomID)
	}

	// Broadcast ReloadRoomsDetail to all clients
	roomIDBytes, err := json.Marshal(event.RoomID)
	if err != nil {
		log.Printf("Error marshalling room ID for ReloadRoomsDetail: %v", err)
	} else {
		go func() {
			h.Broadcast <- NewMessageEvent(
				MessageTypeGlobal,
				NewMessageEventPayload(ReloadRoomsDetail, json.RawMessage(roomIDBytes)),
			)
		}()
	}
}

// findClient 安全地从 hub 的 Users map 中检索客户端。
func (h *Hub) findClient(userId, deviceId string) (*Client, error) {
	h.mu.RLock()         // 获取读锁
	defer h.mu.RUnlock() // 释放读锁

	return h.findClientNoLock(userId, deviceId)
}

// findClientNoLock 是一个在不获取锁的情况下查找客户端的辅助函数。
// 它假定调用者已经持有了适当的 hub 锁。
func (h *Hub) findClientNoLock(userId, deviceId string) (*Client, error) {
	if userClients, ok := h.Users[userId]; ok {
		for c := range userClients {
			if c.DeviceId == deviceId {
				return c, nil
			}
		}
	}
	return nil, ClientNotFindError
}
