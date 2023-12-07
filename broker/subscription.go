package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/meeron/mebrox/logger"
)

const (
	DefaultMaxConsumed        = 3
	DefaultLockTimeoutMinutes = 5
)

type Subscription struct {
	cfg        *subscriptionCfg
	messages   []*Message
	deadLetter []*Message
	mux        *sync.Mutex
}

type subscriptionCfg struct {
	maxConsumed int
	lockTimeout time.Duration
}

func NewSubscription() *Subscription {
	cfg := &subscriptionCfg{
		maxConsumed: DefaultMaxConsumed,
		lockTimeout: DefaultLockTimeoutMinutes * time.Minute,
	}

	return &Subscription{
		cfg:        cfg,
		messages:   make([]*Message, 0),
		deadLetter: make([]*Message, 0),
		mux:        new(sync.Mutex),
	}
}

func (s *Subscription) AddMessage(msg *Message) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.messages = append(s.messages, msg)
}

func (s *Subscription) CommitMessage(id string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	index := len(s.messages)
	for i := 0; i < len(s.messages); i++ {
		if s.messages[i].Id == id {
			index = i
			break
		}
	}

	// If not found sort.Search returns n=len(s.messages)
	if index >= len(s.messages) {
		fmt.Printf("%v %v %v", index, len(s.messages), id)
		return false
	}

	s.messages = append(s.messages[:index], s.messages[index+1:]...)
	logger.Debug("Message committed (%s)", id)

	return true
}

func (s *Subscription) Subscribe(ctx context.Context) <-chan *Message {
	msgChan := make(chan *Message)

	go func(sub *Subscription, c chan *Message, ctx context.Context) {
		defer close(c)

		for {
			for i, msg := range s.messages {
				if msg.lockTime != (time.Time{}) &&
					time.Since(msg.lockTime) < sub.cfg.lockTimeout {
					continue
				}

				if msg.consumedCount >= sub.cfg.maxConsumed {
					sub.mux.Lock()

					// Remove message
					sub.messages = append(sub.messages[:i], sub.messages[i+1:]...)

					// Append message to dead letter
					sub.deadLetter = append(sub.deadLetter, msg)
					sub.mux.Unlock()

					logger.Debug("Message moved to dead letter (%s)", msg.Id)
					continue
				}

				c <- msg
				msg.lockTime = time.Now()
				msg.consumedCount++
			}

			time.Sleep(250 * time.Millisecond)
			if err := ctx.Err(); err != nil {
				return
			}
		}
	}(s, msgChan, ctx)

	return msgChan
}

func (s *Subscription) getMessage(index int) *Message {
	s.mux.Lock()
	defer s.mux.Unlock()

	if index >= len(s.messages) {
		return nil
	}

	return s.messages[index]
}
