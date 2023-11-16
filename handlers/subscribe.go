package handlers

import (
	"fmt"
	"net/http"

	"github.com/meeron/mebrox/server"
)

func Subscribe(s *server.Server, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")

	id := s.Subscribe(w)
	defer s.Unsubscribe(id)
	fmt.Printf("got new connection: %v\n", id)

	s.SendEventTo(id, "ack", id)

	for {
		select {
		case <-r.Context().Done():
			fmt.Printf("connection closed: %v\n", id)
			return nil
		}
	}
}
