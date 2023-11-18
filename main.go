package main

import (
	"github.com/meeron/mebrox/broker"
	"github.com/meeron/mebrox/handlers"
	"github.com/meeron/mebrox/server"
)

func main() {
	broker := broker.NewBroker()
	broker.CreateTopic("test")
	broker.CreateSubscription("test", "test")

	s := server.New(broker)

	s.HandleFunc("/publish", handlers.Publish)
	s.HandleFunc("/subscribe", handlers.Subscribe)

	s.Run(":3000")
}
