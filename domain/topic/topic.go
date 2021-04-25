package topic

import (
	"errors"

	"github.com/gorilla/websocket"
)

// imports

type Topic struct {
	Name        string
	Subscribers []*websocket.Conn
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
