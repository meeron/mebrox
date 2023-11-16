package main

import (
	"github.com/meeron/mebrox/handlers"
	"github.com/meeron/mebrox/server"
)

func main() {
	s := server.New()

	s.HandleFunc("/publish", handlers.Publish)
	s.HandleFunc("/subscribe", handlers.Subscribe)

	s.Run(":3000")
}
