package types

import "encoding/json"

type PayloadType int

const (
	PayloadTypeText PayloadType = iota
	PayloadTypeImage
	PayloadTypeFile
)

type Payload struct {
	Type    PayloadType     `json:"type"`
	Content json.RawMessage `json:"data"`
	File    []byte          `json:"file"`
}

func NewPayload(t PayloadType, content []byte) *Payload {
	return &Payload{
		Type:    t,
		Content: content,
	}
}
