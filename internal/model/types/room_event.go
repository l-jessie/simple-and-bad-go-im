package types

type CreateRoomEvent struct {
	UserID string `json:"userId"`
	RoomID string `json:"roomId"`
	Room   *Room  `json:"room"`
}

type JoinRoomEvent struct {
	RoomID   string  `json:"roomId"`
	UserID   string  `json:"userId"`
	DeviceID string  `json:"deviceId"`
	Client   *Client `json:"client"`
}

type UnJoinRoomEvent struct {
	RoomID   string `json:"roomId"`
	UserID   string `json:"userId"`
	DeviceID string `json:"deviceId"`
}
