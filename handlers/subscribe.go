package handlers

import (
	"net/http"

	"github.com/meeron/mebrox/logger"
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

	s.SendEvent(w, server.Event{
		EventType: "welcome",
	})
	logger.Debug("Connected (%s-%s)", topic, subscription)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")

	for {
		select {
		case <-r.Context().Done():
			logger.Debug("Disonnected (%s-%s)", topic, subscription)
			return nil
		}
	}
}
