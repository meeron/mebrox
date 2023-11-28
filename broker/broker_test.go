package broker

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCreateTopic(t *testing.T) {
	broker := NewBroker()

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
	broker := NewBroker()
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

func TestSubscribe(t *testing.T) {
	broker := NewBroker()

	if err := broker.CreateTopic("test"); err != nil {
		panic(err)
	}
	if err := broker.CreateSubscription("test", "test"); err != nil {
		panic(err)
	}

	t.Run("should subscribe to existing subscription", func(t *testing.T) {
		sub, err := broker.Subscribe("test", "test")

		assert.Nilf(t, err, "%v", err)
		assert.NotNil(t, sub, "subscription IS nil")
	})

	t.Run("should return error when subscribing to not existing", func(t *testing.T) {
		_, err := broker.Subscribe("test", "test1")

		assert.NotNil(t, err, "error SHOULD be nil")
	})
}

func TestBroker_FindSubscription(t *testing.T) {
	broker := NewBroker()

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
	broker := NewBroker()

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
