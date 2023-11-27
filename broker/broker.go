package broker

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/meeron/mebrox/logger"
)

const (
	DefaultMaxConsumed        = 3
	DefaultLockTimeoutMinutes = 5
)

type Broker struct {
	topics map[string]*Topic
}

type Message struct {
	Id            string
	Body          []byte
	lockTime      time.Time
	consumedCount int
}

type Subscription struct {
	cfg        *subscriptionCfg
	messages   []*Message
	deadLetter []*Message
	Msg        chan *Message
}

type Topic struct {
	subscriptions map[string]*Subscription
}

type subscriptionCfg struct {
	maxConsumed int
	lockTimeout time.Duration
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
		topics: make(map[string]*Topic),
	}
}

func (b *Broker) SendMessage(topic string, msg *Message) error {
	t, ok := b.topics[topic]
	if !ok {
		return errors.New("topic not found")
	}

	for _, sub := range t.subscriptions {
		sub.messages = append(sub.messages, msg)
	}
	return nil
}

func (b *Broker) CreateTopic(name string) error {
	_, ok := b.topics[name]
	if ok {
		return errors.New("topic already exists")
	}

	b.topics[name] = &Topic{
		subscriptions: make(map[string]*Subscription),
	}

	return nil
}

func (b *Broker) CreateSubscription(topic string, sub string) error {
	cfg := &subscriptionCfg{
		maxConsumed: DefaultMaxConsumed,
		lockTimeout: DefaultLockTimeoutMinutes * time.Minute,
	}

	t, ok := b.topics[topic]
	if !ok {
		return errors.New("topic does not exists")
	}

	_, ok = t.subscriptions[sub]
	if ok {
		return errors.New("subscription already exists")
	}

	t.subscriptions[sub] = &Subscription{
		cfg:        cfg,
		messages:   make([]*Message, 0),
		deadLetter: make([]*Message, 0),
		Msg:        make(chan *Message),
	}
	go monitor(t.subscriptions[sub])

	return nil
}

func (b *Broker) GetSubscription(topic string, subscription string) (*Subscription, error) {
	t, ok := b.topics[topic]
	if !ok {
		return nil, errors.New("topic not found")
	}

	sub, ok := t.subscriptions[subscription]
	if !ok {
		return nil, errors.New("subscription does not exists")
	}

	return sub, nil
}

func (b *Broker) CommitMessage(topic string, sub string, id string) (bool, error) {
	t, ok := b.topics[topic]
	if !ok {
		return false, errors.New("topic not found")
	}

	s, subOk := t.subscriptions[sub]
	if !subOk {
		return false, errors.New("subscription not found")
	}

	for index, msg := range s.messages {
		if msg.Id == id {
			s.messages = append(s.messages[:index], s.messages[index+1:]...)
			logger.Debug("Message commited (%s)", id)
			return true, nil
		}
	}

	return false, nil
}

func monitor(sub *Subscription) {
	for {
		for i, msg := range sub.messages {
			if msg.lockTime != (time.Time{}) &&
				time.Since(msg.lockTime) < sub.cfg.lockTimeout {
				continue
			}

			if msg.consumedCount >= sub.cfg.maxConsumed {
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
