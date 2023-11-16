package handlers

import (
	"fmt"
	"net/http"
)

func responseBadRequest(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, message)
	return nil
}
