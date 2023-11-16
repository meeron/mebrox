package store

type memoryStore struct{}

func (m *memoryStore) CreateTopic(name string) error {
	return nil
}

func (m *memoryStore) DeleteTopic(name string) error {
	return nil
}

func (m *memoryStore) CreateSubscription(name string) error {
	return nil
}

func (m *memoryStore) DeleteSubscription(name string) error {
	return nil
}

func (m *memoryStore) SaveMessage(topic string, msg Message) error {
	return nil
}

func (m *memoryStore) GetMessages(topic string, subscription string) ([]Message, error) {
	return make([]Message, 0), nil
}
