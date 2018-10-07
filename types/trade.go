package types

import "encoding/json"

type Trade struct {
	MakerOrderID string `json:"maker_order_id"`
	Price        uint64 `json:"price"`
}

func (trade *Trade) FromJson(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

func (trade *Trade) ToJson() []byte {
	str, _ := json.Marshal(trade)
	return str
}
