package handlers

import (
	"errors"
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
	commitMessageRegex  *regexp.Regexp
	topicRegex          *regexp.Regexp
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

	commitParams := commitMessageRegex.FindStringSubmatch(r.URL.String())
	if len(commitParams) > 1 {
		return commitMessage(s, w, r, commitParams)
	}

	topicParams := topicRegex.FindStringSubmatch(r.URL.String())
	if len(topicParams) > 1 {
		return topicCreate(s, w, r, topicParams)
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

	return send(w, msg.Id)
}

func subscribe(s *server.Server, w http.ResponseWriter, r *http.Request, params []string) error {
	if r.Method != http.MethodGet {
		return responseMethodNotAllowed(w)
	}

	topic := params[1]
	subscription := params[2]

	sub, err := s.Broker().Subscribe(topic, subscription)
	if err != nil {
		return err
	}

	if err := s.SendEvent(w, server.Event{
		EventType: "welcome",
	}); err != nil {
		return err
	}
	logger.Debug("Connected (%s-%s)", topic, subscription)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")

	for {
		select {
		case <-r.Context().Done():
			logger.Debug("Disconnected (%s-%s)", topic, subscription)
			return nil
		case msg := <-sub.Msg:
			_ = s.SendEvent(w, server.Event{
				Id:        msg.Id,
				EventType: "message",
				Data:      msg.Body,
			})
			break
		}
	}
}

func commitMessage(s *server.Server, w http.ResponseWriter, r *http.Request, params []string) error {
	if r.Method != http.MethodPost {
		return responseMethodNotAllowed(w)
	}

	topic := params[1]
	subName := params[2]
	id := params[3]

	sub := s.Broker().FindSubscription(topic, subName)
	if sub == nil {
		return errors.New("subscription not found")
	}

	if ok := sub.CommitMessage(id); !ok {
		return responseNotFound(w, "message not found")
	}

	return send(w, "ok")
}

func topicCreate(s *server.Server, w http.ResponseWriter, r *http.Request, params []string) error {
	name := params[1]

	// Create topic
	if r.Method == http.MethodPost {
		if err := s.Broker().CreateTopic(name); err != nil {
			return err
		}

		return responseCreated(w, "created")
	}

	return responseMethodNotAllowed(w)
}

func init() {
	publishMessageRegex = regexp.MustCompile("/topics/(\\w+)/messages($|/)")
	subscribeRegex = regexp.MustCompile("/topics/(\\w+)/subscriptions/(\\w+)/subscribe($|/)")
	commitMessageRegex = regexp.MustCompile("/topics/(\\w+)/subscriptions/(\\w+)/messages/(\\w+)/commit($|/)")
	topicRegex = regexp.MustCompile("/topics/(\\w+)($|/)")
}
