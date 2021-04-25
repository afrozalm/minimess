package topic

import (
	"errors"

	"github.com/gorilla/websocket"
)

// imports

type Topic struct {
	id          int
	name        string
	subscribers []*websocket.Conn
}

var ErrUnknown = errors.New("unknown topic")

type Repository interface {
	// GetAll returns all topics added by the subscribers
	GetAll() []Topic
	// Get topic with the given name
	Get(string) (Topic, error)
	// Add adds a topic to the repository of topics
	Add(Topic) error
}
