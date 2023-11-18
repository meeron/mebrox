package broker

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/meeron/mebrox/logger"
)

const (
	MaxConsumed        = 3
	LockTimeoutMinutes = 1
)

type Broker struct {
	data map[string]map[string]*Subscription
}

type Message struct {
	Id            string
	Body          []byte
	lockTime      time.Time
	consumedCount int
}

type Subscription struct {
	messages   []*Message
	deadLetter []*Message
	Msg        chan *Message
}

func NewMessage(body []byte) *Message {
	return &Message{
		Id:       newMessageId(),
		Body:     body,
		lockTime: time.Time{},
	}
}

func NewBroker() *Broker {
	return &Broker{
		data: make(map[string]map[string]*Subscription),
	}
}

func (b *Broker) SendMessage(topic string, msg *Message) error {
	for _, sub := range b.data[topic] {
		sub.messages = append(sub.messages, msg)
	}
	return nil
}

func (b *Broker) CreateTopic(name string) {
	b.data[name] = make(map[string]*Subscription)
}

func (b *Broker) CreateSubscription(topic string, sub string) {
	b.data[topic][sub] = &Subscription{
		messages:   make([]*Message, 0),
		deadLetter: make([]*Message, 0),
		Msg:        make(chan *Message),
	}
	go monitor(b.data[topic][sub])
}

func (b *Broker) GetSubscription(topic string, subscription string) *Subscription {
	t, topicOk := b.data[topic]
	if !topicOk {
		return nil
	}

	sub, ok := t[subscription]
	if !ok {
		return nil
	}

	return sub
}

func monitor(sub *Subscription) {
	for {
		for i, msg := range sub.messages {
			if msg.lockTime != (time.Time{}) &&
				time.Since(msg.lockTime).Minutes() < LockTimeoutMinutes {
				continue
			}

			if msg.consumedCount >= MaxConsumed {
				sub.messages = append(sub.messages[:i], sub.messages[i+1:]...)
				sub.deadLetter = append(sub.deadLetter, msg)

				logger.Debug("Message moved to dead letter (%s)", msg.Id)
				continue
			}

			sub.Msg <- msg

			msg.lockTime = time.Now()
			msg.consumedCount++
		}

		time.Sleep(1 * time.Second)
	}
}

func newMessageId() string {
	data := make([]byte, 16)

	rand.Read(data)

	return hex.EncodeToString(data)
}
