package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"net/http"
	"order_engine/engine"
	"order_engine/rabbitmq"
	"order_engine/types"
	"sync"
)

const (
	TradeChannel = "trades"
	OrderChannel = "orders"
)

var count int
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Conn struct {
	*websocket.Conn
	mu sync.Mutex
}

var conn *Conn
var book *engine.OrderBook

// MessagingClient instance
var MessagingClient rabbitmq.IMessagingClient

// ConnectionEndpoint is the the handleFunc function for websocket connections
// It handles incoming websocket messages and routes the message according to
// channel parameter in channelMessage
func ConnectionEndpoint(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn = &Conn{c, sync.Mutex{}}
	count++
	fmt.Println(count, r.Host, r.Body)
	go func() {
		// Recover in case of any panic in websocket. So that the app doesn't crash ===
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if err != nil {
					return
				}

				if !ok {
				}
			}
		}()

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
			}
			if messageType != 1 {
				return
			}
			MessagingClient.PublishOnQueue(p, "orders")
		}
	}()
}

func NewConnection(conn *websocket.Conn) *Conn {
	return &Conn{conn, sync.Mutex{}}
}

// SendMessage constructs the message with proper structure to be sent over websocket
func SendMessage(message []byte) {
	var order types.Order
	// decode the message
	order.FromJson(message)

	// process the order
	trades := book.Process(&order)
	// send trades to message queue
	for _, trade := range trades {
		rawTrade := trade.ToJson()
		fmt.Println("raw trade", string(rawTrade))
		conn.mu.Lock()
		defer conn.mu.Unlock()
		err := conn.WriteJSON(message)
		if err != nil {
			conn.Close()
		}
	}
}

func InitializeMessaging(ob *engine.OrderBook) {
	book = ob
	MessagingClient = &rabbitmq.MessagingClient{}
	MessagingClient.ConnectToBroker("amqp://guest:guest@localhost:5672/")
	MessagingClient.SubscribeToQueue("orders", "orderconsumer", onMessage)
}

func onMessage(delivery amqp.Delivery) {
	SendMessage(delivery.Body)
}
