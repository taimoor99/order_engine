package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"order_engine/engine"
)

var KafkaConsumer *cluster.Consumer
var KafkaProducer sarama.AsyncProducer

func NewKafkaCP(done chan bool) {

	// create the order book
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}

	// create the consumer and listen for new order messages
	KafkaConsumer = CreateConsumer()

	// create the producer of trade messages
	KafkaProducer = CreateProducer()

	// start processing orders
	go func() {
		for msg := range KafkaConsumer.Messages() {
			fmt.Println("message", msg.Topic, string(msg.Value))
			var order engine.Order
			// decode the message
			order.FromJson(msg.Value)
			// process the order
			trades := book.Process(order)
			// send trades to message queue
			for _, trade := range trades {
				rawTrade := trade.ToJson()
				fmt.Println("raw trade", string(rawTrade))
				KafkaProducer.Input() <- &sarama.ProducerMessage{
					Topic: "trades",
					Value: sarama.ByteEncoder("one order match"),
				}
			}
			// mark the message as processed
			KafkaConsumer.MarkOffset(msg, "")
		}
		done <- true
	}()

}

func ProduceMsg(msg []byte) {
	KafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.ByteEncoder(msg),
	}
}

//
// Create the consumer
//

func CreateConsumer() *cluster.Consumer {
	// define our configuration to the cluster
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = false
	config.Group.Return.Notifications = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// create the consumer
	consumer, err := cluster.NewConsumer([]string{"127.0.0.1:9092"}, "myconsumer", []string{"orders"}, config)
	if err != nil {
		log.Fatal("Unable to connect consumer to kafka cluster")
	}
	go handleErrors(consumer)
	go handleNotifications(consumer)
	return consumer

}

func handleErrors(consumer *cluster.Consumer) {
	for err := range consumer.Errors() {
		log.Printf("Error: %s\n", err.Error())
	}
}

func handleNotifications(consumer *cluster.Consumer) {
	for ntf := range consumer.Notifications() {
		log.Printf("Rebalanced: %+v\n", ntf)
	}
}

//
// Create the producer
//

func CreateProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	producer, err := sarama.NewAsyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatal("Unable to connect producer to kafka server")
	}
	return producer
}
