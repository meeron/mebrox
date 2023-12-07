package main

import (
	"context"
	"github.com/meeron/mebrox/broker"
	"github.com/meeron/mebrox/handlers"
	"github.com/meeron/mebrox/server"
)

func main() {
	br := broker.NewBroker(context.Background())
	s := server.New(br)

	s.HandleFunc("/topics", handlers.Topics)
	s.HandleFunc("/topics/", handlers.Topics)

	s.Run(":3000")
}
