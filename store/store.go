package store

import (
	"crypto/rand"
	"encoding/hex"
)

type Storer interface {
	CreateTopic(name string) error
	DeleteTopic(name string) error
	CreateSubscription(topic string, subscription string) error
	DeleteSubscription(topic string, subscription string) error
	SaveMessage(topic string, msg Message) error
	GetMessages(topic string, subscription string) ([]Message, error)
}

type Message struct {
	Id      string
	Headers map[string]string
	Body    []byte
}

func NewMessage(body []byte) Message {
	return Message{
		Id:      newMessageId(),
		Headers: make(map[string]string),
		Body:    body,
	}
}

func New() Storer {
	return &memoryStore{
		data: make(map[string]map[string][]Message),
	}
}

func newMessageId() string {
	data := make([]byte, 16)

	rand.Read(data)

	return hex.EncodeToString(data)
}
