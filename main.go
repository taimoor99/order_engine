package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"order_engine/engine"
	"order_engine/redis"
	"order_engine/ws"
)

func main() {

	redicon := redis.NewRedisConnection("redis://localhost:6379")
	// create the order book
	book := &engine.OrderBook{
		RedisConn: redicon,
	}

	ws.InitializeMessaging(book)

	router := NewRouter()
	// http.Handle("/", router)
	router.HandleFunc("/socket", ws.ConnectionEndpoint)

	// start the server
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	panic(http.ListenAndServe(":9090", handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router)))

}

func NewRouter() *mux.Router {

	r := mux.NewRouter()

	// cronService.InitCrons()
	return r
}
