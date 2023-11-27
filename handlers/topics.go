package handlers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/meeron/mebrox/broker"
	"github.com/meeron/mebrox/logger"
	"github.com/meeron/mebrox/server"
)

var (
	publishMessageRegex *regexp.Regexp
	subscribeRegex      *regexp.Regexp
)

func Topics(s *server.Server, w http.ResponseWriter, r *http.Request) error {
	publishParams := publishMessageRegex.FindStringSubmatch(r.URL.String())
	if len(publishParams) > 1 {
		return publishMessage(s, w, r, publishParams)
	}

	subscribeParams := subscribeRegex.FindStringSubmatch(r.URL.String())
	if len(subscribeParams) > 1 {
		return subscribe(s, w, r, subscribeParams)
	}

	return responseNotFound(w, "not found")
}

func publishMessage(s *server.Server, w http.ResponseWriter, r *http.Request, params []string) error {
	if r.Method != http.MethodPost {
		return responseMethodNotAllowed(w)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	topic := params[1]

	msg := broker.NewMessage(body)
	if err := s.Broker().SendMessage(topic, msg); err != nil {
		return err
	}

	fmt.Fprint(w, "ok")
	return nil
}

func subscribe(s *server.Server, w http.ResponseWriter, r *http.Request, params []string) error {
	if r.Method != http.MethodGet {
		return responseMethodNotAllowed(w)
	}

	topic := params[1]
	subscription := params[2]

	sub, err := s.Broker().GetSubscription(topic, subscription)
	if err != nil {
		return err
	}

	s.SendEvent(w, server.Event{
		EventType: "welcome",
	})
	logger.Debug("Connected (%s-%s)", topic, subscription)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")

	for {
		select {
		case <-r.Context().Done():
			logger.Debug("Disonnected (%s-%s)", topic, subscription)
			return nil
		case msg := <-sub.Msg:
			s.SendEvent(w, server.Event{
				Id:        msg.Id,
				EventType: "message",
				Data:      msg.Body,
			})
			break
		}
	}
}

func init() {
	publishMessageRegex = regexp.MustCompile("/topics/(\\w+)/messages($|/)")
	subscribeRegex = regexp.MustCompile("/topics/(\\w+)/subscriptions/(\\w+)/subscribe($|/)")
}
