package types

type RegisterEvent struct {
	UserId string  `json:"userId"`
	Client *Client `json:"client"`
}

type UnRegisterEvent struct {
	UserId string  `json:"userId"`
	Client *Client `json:"client"`
}
