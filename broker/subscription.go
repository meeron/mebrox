package broker

import (
	"sort"
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
	Msg        chan *Message
}

type subscriptionCfg struct {
	maxConsumed int
	lockTimeout time.Duration
}

func (s *Subscription) AddMessage(msg *Message) {
	s.messages = append(s.messages, msg)
}

func (s *Subscription) CommitMessage(id string) bool {
	index := sort.Search(len(s.messages), func(i int) bool {
		return s.messages[i].Id == id
	})

	// If not found sort.Search returns n=len(s.messages)
	if index >= len(s.messages) {
		return false
	}

	s.messages = append(s.messages[:index], s.messages[index+1:]...)
	logger.Debug("Message committed (%s)", id)

	return true
}

func monitor(sub *Subscription) {
	for {
		for index, msg := range sub.messages {
			if msg.lockTime != (time.Time{}) &&
				time.Since(msg.lockTime) < sub.cfg.lockTimeout {
				continue
			}

			if msg.consumedCount >= sub.cfg.maxConsumed {
				// Remove message
				sub.messages = append(sub.messages[:index], sub.messages[index+1:]...)

				// Append message to dead letter
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
