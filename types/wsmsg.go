package types

type WebSocketMessage struct {
	Channel string           `json:"channel"`
	Payload WebSocketPayload `json:"payload"`
}

type WebSocketPayload struct {
	Type string      `json:"type"`
	Hash string      `json:"hash,omitempty"`
	Data interface{} `json:"data"`
}
