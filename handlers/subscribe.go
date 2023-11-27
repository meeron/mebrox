package handlers

import (
	"net/http"

	"github.com/meeron/mebrox/logger"
	"github.com/meeron/mebrox/server"
)

func Subscribe(s *server.Server, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return responseMethodNotAllowed(w)
	}

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		return responseBadRequest(w, "Invalid topic name")
	}

	subscription := r.URL.Query().Get("subscription")
	if subscription == "" {
		return responseBadRequest(w, "Invalid subscription name")
	}

	sub, err := s.Broker().GetSubscription(topic, subscription)
	if err != nil {
		return err
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
		case msg := <-sub.Msg:
			s.SendEvent(w, server.Event{
				Id:        msg.Id,
				EventType: "message",
				Data:      msg.Body,
			})
			break
		}
	}
}
