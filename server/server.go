package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/meeron/mebrox/store"
)

type Server struct {
	mux     *http.ServeMux
	clients map[string]map[string]http.ResponseWriter
	store   store.Storer
}

type event struct {
	id        string
	eventType string
	data      []byte
}

type ServerHandler func(s *Server, w http.ResponseWriter, r *http.Request) error

func New() *Server {
	return &Server{
		clients: make(map[string]map[string]http.ResponseWriter),
		mux:     http.NewServeMux(),
		store:   store.New(),
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

func (s *Server) SendMessage(topic string, body []byte) error {
	msg := store.NewMessage(body)

	if err := s.store.EnsureTopic(topic); err != nil {
		return err
	}

	if err := s.store.SaveMessage(topic, msg); err != nil {
		return err
	}

	return s.sendToTopic(topic, event{
		id:        msg.Id,
		eventType: "message",
		data:      msg.Body,
	})
}

func (s *Server) Subscribe(w http.ResponseWriter, topic string, subscription string) error {
	if _, ok := s.clients[subscription]; ok {
		return fmt.Errorf("Subscription '%s' already exists", subscription)
	}

	if err := s.store.EnsureSubscription(topic, subscription); err != nil {
		return err
	}

	messages, err := s.store.GetMessages(topic, subscription)
	if err != nil {
		return err
	}

	if _, ok := s.clients[topic]; !ok {
		s.clients[topic] = make(map[string]http.ResponseWriter)
	}

	s.clients[topic][subscription] = w
	log.Default().Printf("[Debug] Subscribed '%s' to '%s'", subscription, topic)

	for _, msg := range messages {
		e := event{
			id:        msg.Id,
			eventType: "message",
			data:      msg.Body,
		}
		if err := s.sendToSubscription(topic, subscription, e); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Unsubscribe(topic string, subscription string) {
	delete(s.clients, subscription)
	log.Default().Printf("[Debug] Unsubscribed '%s' from '%s'", subscription, topic)
}

func (s *Server) sendToSubscription(topic string, subsription string, e event) error {
	subscriptions := s.clients[topic]
	w := subscriptions[subsription]

	fmt.Fprintf(w, "event: %s\n", e.eventType)
	if e.id != "" {
		fmt.Fprintf(w, "id: %s\n", e.id)
	}
	fmt.Fprintf(w, "data: %s", e.data)
	fmt.Fprint(w, "\n\n")

	f, ok := w.(http.Flusher)
	if ok {
		f.Flush()
	}

	return nil
}

func (s *Server) sendToTopic(topic string, e event) error {
	for subscription := range s.clients[topic] {
		if err := s.sendToSubscription(topic, subscription, e); err != nil {
			return err
		}
	}

	return nil
}
