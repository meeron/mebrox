package main

import (
	"log"

	"github.com/meeron/mebrox/broker"
	"github.com/meeron/mebrox/handlers"
	"github.com/meeron/mebrox/server"
)

func main() {
	broker := broker.NewBroker()
	if err := broker.CreateTopic("test"); err != nil {
		log.Fatal(err)
	}
	if err := broker.CreateSubscription("test", "test"); err != nil {
		log.Fatal(err)
	}

	s := server.New(broker)

	s.HandleFunc("/publish", handlers.Publish)
	s.HandleFunc("/subscribe", handlers.Subscribe)
	s.HandleFunc("/messages", handlers.Messages)
	s.HandleFunc("/messages/", handlers.Messages)

	s.Run(":3000")
}
