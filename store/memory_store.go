package store

import (
	"fmt"
)

type memoryStore struct {
	data map[string]map[string][]Message
}

func (m *memoryStore) EnsureTopic(name string) error {
	if _, ok := m.data[name]; ok {
		return nil
	}

	m.data[name] = make(map[string][]Message)
	return nil
}

func (m *memoryStore) DeleteTopic(name string) error {
	delete(m.data, name)
	return nil
}

func (m *memoryStore) EnsureSubscription(topic string, subscription string) error {
	if err := m.EnsureTopic(topic); err != nil {
		return err
	}

	subscriptions := m.data[topic]

	if _, ok := subscriptions[subscription]; ok {
		return nil
	}

	subscriptions[subscription] = make([]Message, 0)

	return nil
}

func (m *memoryStore) DeleteSubscription(topic string, subscription string) error {
	subscriptions, ok := m.data[topic]
	if !ok {
		return nil
	}

	delete(subscriptions, subscription)

	return nil
}

func (m *memoryStore) SaveMessage(topic string, msg Message) error {
	subscriptions, ok := m.data[topic]
	if !ok {
		return fmt.Errorf("topic '%s' not found", topic)
	}

	for name := range subscriptions {
		subscriptions[name] = append(subscriptions[name], msg)
	}

	return nil
}

func (m *memoryStore) GetMessages(topic string, subscription string) ([]Message, error) {
	subscriptions, ok := m.data[topic]
	if !ok {
		return make([]Message, 0), fmt.Errorf("topic '%s' not found", topic)
	}

	messages, ok := subscriptions[subscription]
	if !ok {
		return make([]Message, 0), nil
	}

	return messages, nil
}
