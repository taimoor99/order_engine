package types

import (
	"encoding/json"
	"fmt"
)

type Order struct {
	Amount  uint64 `json:"amount"`
	Price   uint64 `json:"price"`
	BaseID  string `json:"baseid"`
	TokenID string `json:"tokenid"`
	Side    string `json:"side"`
}

func (order *Order) FromJson(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJson() []byte {
	str, _ := json.Marshal(order)
	return str
}

// GetOBKeys returns the keys corresponding to an order
// orderbook price point key
// orderbook list key corresponding to order price.
func (o *Order) GetOBKeys() (ss string) {
	var k string
	fmt.Println("order type", o.Side)
	if o.Side == "BUY" {
		k = "SELL"
	} else if o.Side == "SELL" {
		k = "BUY"
	}

	ss = "*" + "::orders::" + k + "*"
	return
}

// GetOBMatchKey returns the orderbook price point key
// aginst which the order needs to be matched
func (o *Order) GetOBMatchKey() (ss string) {
	ss = o.GetKVPrefix() + "::orders::" + o.Side
	return
}

// GetKVPrefix returns the key value store(redis) prefix to be used
// by matching engine correspondind to a particular order.
func (o *Order) GetKVPrefix() string {
	return o.BaseID + "::" + o.TokenID
}
