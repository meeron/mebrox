package broker

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Broker struct {
	topics        map[string]*Topic
	subscriptions map[string]*Subscription
	subMux        sync.Mutex
	createMux     sync.Mutex
}

type Message struct {
	Id            string
	Body          []byte
	lockTime      time.Time
	consumedCount int
}

type Topic struct {
	subscriptions map[string]*Subscription
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
		topics:        make(map[string]*Topic),
		subscriptions: make(map[string]*Subscription),
	}
}

func (b *Broker) SendMessage(topic string, msg *Message) error {
	t, ok := b.topics[topic]
	if !ok {
		return errors.New("topic not found")
	}

	for _, sub := range t.subscriptions {
		sub.AddMessage(msg)
	}
	return nil
}

func (b *Broker) CreateTopic(name string) error {
	b.createMux.Lock()
	defer b.createMux.Unlock()

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
	b.createMux.Lock()
	defer b.createMux.Unlock()

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

	return nil
}

func (b *Broker) Subscribe(topic string, subscription string) (*Subscription, error) {
	t, ok := b.topics[topic]
	if !ok {
		return nil, errors.New("topic not found")
	}

	sub, ok := t.subscriptions[subscription]
	if !ok {
		return nil, errors.New("subscription does not exists")
	}

	key := fmt.Sprintf("%s_%s", topic, subscription)

	b.subMux.Lock()
	_, ok = b.subscriptions[key]
	if !ok {
		b.subscriptions[key] = sub
		go monitor(sub)
	}
	b.subMux.Unlock()

	return sub, nil
}

func (b *Broker) FindSubscription(topic string, sub string) *Subscription {
	t, ok := b.topics[topic]
	if !ok {
		return nil
	}

	s, subOk := t.subscriptions[sub]
	if !subOk {
		return nil
	}

	return s
}

func newMessageId() string {
	data := make([]byte, 16)

	rand.Read(data)

	return hex.EncodeToString(data)
}
