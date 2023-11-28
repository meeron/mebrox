package main

import (
	"github.com/meeron/mebrox/broker"
	"github.com/meeron/mebrox/handlers"
	"github.com/meeron/mebrox/server"
	"log"
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

	s.HandleFunc("/topics", handlers.Topics)
	s.HandleFunc("/topics/", handlers.Topics)

	s.Run(":3000")
}
