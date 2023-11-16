package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	mux     *http.ServeMux
	clients map[string]http.ResponseWriter
}

type ServerHandler func(s *Server, w http.ResponseWriter, r *http.Request) error

func New() *Server {
	return &Server{
		clients: make(map[string]http.ResponseWriter),
		mux:     http.NewServeMux(),
	}
}

func (s *Server) Run(addr string) {
	ser := http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	log.Default().Println("Listening...")
	if err := ser.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) HandleFunc(pattern string, handler ServerHandler) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(s, w, r); err != nil {
			log.Default().Printf("Error %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
		}
	})
}

func (s *Server) SendEvent(event string, data string) error {
	for id := range s.clients {
		if err := s.SendEventTo(id, event, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) SendEventTo(clientId string, event string, data string) error {
	w, ok := s.clients[clientId]
	if !ok {
		return errors.New("client not found")
	}

	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s", data)
	fmt.Fprint(w, "\n\n")

	f, ok := w.(http.Flusher)
	if ok {
		f.Flush()
	}

	return nil
}

func (s *Server) Subscribe(w http.ResponseWriter) string {
	id := newId()
	s.clients[id] = w

	return id
}

func (s *Server) Unsubscribe(id string) {
	delete(s.clients, id)
}

func newId() string {
	data := make([]byte, 16)

	rand.Read(data)

	return hex.EncodeToString(data)
}
