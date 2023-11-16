package server

import (
	"log"
	"net/http"
)

type Server struct {
	mux     *http.ServeMux
	clients map[string]http.ResponseWriter
}

type ServerHandler func(s *Server, w http.ResponseWriter, r *http.Request)

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
		handler(s, w, r)
	})
}
