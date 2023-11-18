package handlers

import (
	"fmt"
	"io"
	"net/http"

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

	_, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// TODO: Enqueue message

	fmt.Fprint(w, "ok")

	return nil
}
