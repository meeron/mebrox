package handlers

import (
	"fmt"
	"net/http"

	"github.com/meeron/mebrox/server"
)

func Subscribe(s *server.Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//w.Header().Set("Content-Type", "text/event-stream")
	//w.Header().Set("Cache-Control", "no-store")

	//id := s.Subscribe(w)
	//fmt.Printf("got new connection: %v\n", id)

	//s.SendEventTo(id.String(), "ack", id.String())

	/*
		for {
			select {
			case <-r.Context().Done():
				s.Unsubscribe(id)
				fmt.Printf("connection closed: %v\n", id)
				return
			}
		}
	*/

	fmt.Fprint(w, "ok")
}
