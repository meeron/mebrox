package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/meeron/mebrox/broker"
	"github.com/meeron/mebrox/server"
)

func Publish(s *server.Server, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		return responseBadRequest(w, "Invalid topic name")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	msg := broker.NewMessage(body)
	if err := s.Broker().SendMessage(topic, msg); err != nil {
		return err
	}

	fmt.Fprint(w, "ok")

	return nil
}
