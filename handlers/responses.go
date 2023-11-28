package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func responseNotFound(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusNotFound)
	return send(w, message)
}

func responseMethodNotAllowed(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

func responseCreated(w http.ResponseWriter, msg string) error {
	w.WriteHeader(http.StatusCreated)
	return send(w, msg)
}

func send(w io.Writer, a ...any) error {
	_, err := fmt.Fprint(w, a...)
	return err
}

func sendf(w io.Writer, format string, a ...any) error {
	_, err := fmt.Fprintf(w, format, a...)
	return err
}
