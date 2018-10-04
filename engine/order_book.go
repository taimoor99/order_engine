package engine

import (
	"fmt"
	"order_engine/rabbitmq"
	"order_engine/redis"
	"order_engine/types"
	"sync"
)

// OrderBook type
type OrderBook struct {
	RedisConn    *redis.RedisConnection
	RabbitMQConn *rabbitmq.MessagingClient
	mutex        *sync.Mutex
}

// Add a buy order to the order book
func (book *OrderBook) addOrder(order *types.Order) {
	pricePointSetKey := order.GetOBMatchKey()

	err :=  book.AddToPricePointSet(pricePointSetKey, string(order.Price))
	if err != nil {
		return
	}
	fmt.Println("order added")
}


