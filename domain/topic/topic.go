package topic

import (
	"errors"

	"github.com/gorilla/websocket"
)

type Topic struct {
	Name          string
	SubscriberMap map[string]*websocket.Conn
}

var ErrUnknown = errors.New("unknown topic")
var ErrDuplicate = errors.New("topic already exists")

type Repository interface {
	// GetAll returns all topics added by the subscribers
	GetAll() []*Topic
	// Get topic with the given name
	Get(string) (*Topic, error)
	// Add adds a topic to the repository of topics
	Add(*Topic) error
}

func (t *Topic) AddUserConnectionById(uid string, conn *websocket.Conn) {
	t.SubscriberMap[uid] = conn
}

func (t *Topic) RemoveUserConnectionById(uid string) {
	if _, ok := t.SubscriberMap[uid]; ok {
		delete(t.SubscriberMap, uid)
	}
}

func (t *Topic) GetUserConnectionById(uid string) *websocket.Conn {
	if conn, ok := t.SubscriberMap[uid]; ok {
		return conn
	}
	return nil
}

func (t *Topic) GetAllConnections() []*websocket.Conn {
	res := make([]*websocket.Conn, 5)
	for _, conn := range t.SubscriberMap {
		res = append(res, conn)
	}
	return res
}
