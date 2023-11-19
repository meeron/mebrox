package broker

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/meeron/mebrox/logger"
)

const (
	DefaultMaxConsumed        = 3
	DefaultLockTimeoutMinutes = 5
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
	cfg        *subscriptionCfg
	messages   []*Message
	deadLetter []*Message
	Msg        chan *Message
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
	cfg := &subscriptionCfg{
		maxConsumed: DefaultMaxConsumed,
		lockTimeout: DefaultLockTimeoutMinutes * time.Minute,
	}

	b.data[topic][sub] = &Subscription{
		cfg:        cfg,
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

func (b *Broker) CommitMessage(id string) (bool, error) {
	for _, topic := range b.data {
		for _, sub := range topic {
			for index, msg := range sub.messages {
				if msg.Id == id {
					sub.messages = append(sub.messages[:index], sub.messages[index+1:]...)
					logger.Debug("Message commited (%s)", id)
					return true, nil
				}
			}
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
