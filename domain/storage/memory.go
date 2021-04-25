package storage

import "github.com/afrozalm/minimess/domain/topic"

type MemoryTopicStorage struct {
	// duplication done here for O(1) duplicacy check and O(1) get all topics
	topicMap  map[string]*topic.Topic
	topicList []*topic.Topic
}

func (m *MemoryTopicStorage) Add(t *topic.Topic) error {
	if _, ok := m.topicMap[t.Name]; !ok {
		return topic.ErrDuplicate
	}

	m.topicMap[t.Name] = t
	m.topicList = append(m.topicList, t)

	return nil
}

func (m *MemoryTopicStorage) Get(name string) (*topic.Topic, error) {
	if t, ok := m.topicMap[name]; ok {
		return t, nil
	}
	return nil, topic.ErrUnknown
}

func (m *MemoryTopicStorage) GetAll() []*topic.Topic {
	return m.topicList
}
