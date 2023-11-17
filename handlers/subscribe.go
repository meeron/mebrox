package handlers

import (
	"net/http"

	"github.com/meeron/mebrox/server"
)

func Subscribe(s *server.Server, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		return responseBadRequest(w, "Invalid topic name")
	}

	subscription := r.URL.Query().Get("subscription")
	if subscription == "" {
		return responseBadRequest(w, "Invalid subscription name")
	}

	if err := s.Subscribe(w, topic, subscription); err != nil {
		return err
	}
	defer s.Unsubscribe(topic, subscription)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")

	for {
		select {
		case <-r.Context().Done():
			return nil
		}
	}
}
