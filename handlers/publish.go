package handlers

import (
	"fmt"
	"net/http"

	"github.com/meeron/mebrox/server"
)

func Publish(s *server.Server, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "ok")
}
