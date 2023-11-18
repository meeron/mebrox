package server

import (
	"fmt"
	"net/http"

	"github.com/meeron/mebrox/broker"
	"github.com/meeron/mebrox/logger"
)

type Server struct {
	mux     *http.ServeMux
	clients map[string][]http.ResponseWriter
	broker  *broker.Broker
}

type Event struct {
	Id        string
	EventType string
	Data      []byte
}

type ServerHandler func(s *Server, w http.ResponseWriter, r *http.Request) error

func New(broker *broker.Broker) *Server {
	return &Server{
		clients: make(map[string][]http.ResponseWriter),
		mux:     http.NewServeMux(),
		broker:  broker,
	}
}

func (s *Server) Run(addr string) {
	ser := http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	logger.Info("Listening...")
	if err := ser.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func (s *Server) HandleFunc(pattern string, handler ServerHandler) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(s, w, r); err != nil {
			logger.Error(err)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
		}
	})
}

func (s *Server) SendEvent(w http.ResponseWriter, e Event) error {
	fmt.Fprintf(w, "event: %s\n", e.EventType)
	if e.Id != "" {
		fmt.Fprintf(w, "id: %s\n", e.Id)
	}
	fmt.Fprintf(w, "data: %s", e.Data)
	fmt.Fprint(w, "\n\n")

	f, ok := w.(http.Flusher)
	if ok {
		f.Flush()
	}

	return nil
}

func (s *Server) Broker() *broker.Broker {
	return s.broker
}
