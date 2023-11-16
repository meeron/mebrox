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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := s.SendMessage("", body); err != nil {
		return err
	}

	fmt.Fprint(w, "ok")

	return nil
}
