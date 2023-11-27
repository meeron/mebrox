package broker

import (
	"time"

	"github.com/meeron/mebrox/logger"
)

const (
	DefaultMaxConsumed        = 3
	DefaultLockTimeoutMinutes = 5
)

type Subscription struct {
	cfg        *subscriptionCfg
	messages   map[string]*Message
	deadLetter map[string]*Message
	Msg        chan *Message
}

type subscriptionCfg struct {
	maxConsumed int
	lockTimeout time.Duration
}

func (s *Subscription) AddMessage(msg *Message) {
	s.messages[msg.Id] = msg
}

func (s *Subscription) CommitMessage(id string) bool {
	if _, ok := s.messages[id]; !ok {
		return false
	}

	delete(s.messages, id)

	logger.Debug("Message commited (%s)", id)

	return false
}

func monitor(sub *Subscription) {
	for {
		for _, msg := range sub.messages {
			if msg.lockTime != (time.Time{}) &&
				time.Since(msg.lockTime) < sub.cfg.lockTimeout {
				continue
			}

			if msg.consumedCount >= sub.cfg.maxConsumed {
				delete(sub.messages, msg.Id)
				sub.deadLetter[msg.Id] = msg

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
