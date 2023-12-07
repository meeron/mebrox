package broker

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCreateTopic(t *testing.T) {
	broker := NewBroker(context.TODO())

	t.Run("should create new topic", func(t *testing.T) {
		err := broker.CreateTopic("topic")
		assert.Nilf(t, err, "%v", err)
	})

	t.Run("should return error when topic exists", func(t *testing.T) {
		err := broker.CreateTopic("topic")

		assert.NotNil(t, err, "error should NOT be nil")
	})
}

func TestCreateSubscription(t *testing.T) {
	broker := NewBroker(context.TODO())
	if err := broker.CreateTopic("test"); err != nil {
		panic(err)
	}

	t.Run("should create new subscription", func(t *testing.T) {
		err := broker.CreateSubscription("test", "test")
		assert.Nilf(t, err, "%v", err)
	})

	t.Run("should return error when subscription exists", func(t *testing.T) {
		err := broker.CreateSubscription("test", "test")

		assert.NotNil(t, err, "error should NOT be nil")
	})
}

func TestFindSubscription(t *testing.T) {
	broker := NewBroker(context.TODO())

	if err := broker.CreateTopic("test"); err != nil {
		panic(err)
	}
	if err := broker.CreateSubscription("test", "test"); err != nil {
		panic(err)
	}

	t.Run("should find subscription", func(t *testing.T) {
		sub := broker.FindSubscription("test", "test")

		assert.NotNil(t, sub, "subscription IS nil")
	})
}

func TestBroker_FindSubscription(t *testing.T) {
	broker := NewBroker(context.TODO())

	if err := broker.CreateTopic("test"); err != nil {
		panic(err)
	}
	if err := broker.CreateSubscription("test", "test"); err != nil {
		panic(err)
	}

	t.Run("should find subscription", func(t *testing.T) {
		sub := broker.FindSubscription("test", "test")

		assert.NotNil(t, sub)
	})

	t.Run("should return nil when not exists", func(t *testing.T) {
		sub := broker.FindSubscription("test", "test1")

		assert.Nil(t, sub)
	})
}

func TestAddMessage(t *testing.T) {
	broker := NewBroker(context.TODO())

	if err := broker.CreateTopic("test"); err != nil {
		panic(err)
	}
	if err := broker.CreateSubscription("test", "test"); err != nil {
		panic(err)
	}
	sub := broker.FindSubscription("test", "test")

	t.Run("should add one message", func(t *testing.T) {
		sub.AddMessage(NewMessage([]byte{1}))

		assert.Equal(t, 1, len(sub.messages))

		// Clear messages for next test
		sub.messages = make([]*Message, 0)
	})

	t.Run("should add multiple messages", func(t *testing.T) {
		const messagesToAdd int = 50

		wg := new(sync.WaitGroup)
		wg.Add(messagesToAdd)

		for i := 0; i < messagesToAdd; i++ {
			go func(num int, wg *sync.WaitGroup) {
				sub.AddMessage(NewMessage([]byte{byte(num)}))
				wg.Done()
			}(i, wg)
		}
		wg.Wait()

		assert.Equal(t, messagesToAdd, len(sub.messages))
	})
}

func TestGetMessageAndCommit(t *testing.T) {
	broker := NewBroker(context.TODO())

	if err := broker.CreateTopic("test"); err != nil {
		panic(err)
	}
	if err := broker.CreateSubscription("test", "test"); err != nil {
		panic(err)
	}
	sub := broker.FindSubscription("test", "test")

	t.Run("should add receive messages", func(t *testing.T) {
		const messagesToAdd int = 50
		ctx, cancel := context.WithCancel(context.TODO())

		go func(n int, cancel context.CancelFunc) {
			for i := 0; i < n; i++ {
				sub.AddMessage(NewMessage([]byte{byte(i)}))
			}
		}(messagesToAdd, cancel)

		for msg := range sub.Subscribe(ctx) {
			// Cancel after first read from channel
			cancel()
			if ok := sub.CommitMessage(msg.Id); !ok {
				panic(errors.New("message not committed. " + msg.Id))
			}
		}

		assert.Equal(t, 0, len(sub.messages))
	})
}
