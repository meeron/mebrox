package main

import (
	"github.com/meeron/mebrox/server"
)

func main() {
	s := server.New()

	s.Run(":3000")
}
