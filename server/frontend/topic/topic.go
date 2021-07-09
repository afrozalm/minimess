package topic

import (
	"context"

	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/server/frontend/interfaces"
	"github.com/afrozalm/minimess/sets"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type Topic struct {
	clients      sets.Set
	name         string
	addClient    chan interfaces.Client
	removeClient chan interfaces.Client
	broadcast    chan *message.Chat
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewTopic(name string, parent context.Context) *Topic {
	ctx, cancel := context.WithCancel(parent)
	return &Topic{
		name:         name,
		addClient:    make(chan interfaces.Client, 10),
		removeClient: make(chan interfaces.Client, 10),
		broadcast:    make(chan *message.Chat, 10),
		clients:      make(sets.Set),
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (t Topic) GetName() string {
	return t.name
}

func (t Topic) AddClient(c interfaces.Client) {
	if t.ctx.Err() != nil {
		return
	}
	t.addClient <- c
}

func (t Topic) RemoveClient(c interfaces.Client) {
	if t.ctx.Err() != nil {
		return
	}
	t.removeClient <- c
}

func (t Topic) BroadcastOut(m *message.Chat) {
	if t.ctx.Err() != nil {
		return
	}
	t.broadcast <- m
}

func (t *Topic) Run(closedTopics chan interfaces.Topic) {
	defer func() {
		// garbage collect
		log.Trace("closing topic goroutine for topic ", t.GetName())
		closedTopics <- t
		t.cancel()
		close(t.addClient)
		close(t.broadcast)
		close(t.removeClient)
		log.Trace("closed topic ", t.GetName())
		t.clients = nil
		// make an rpc request to fanout service to unsub host-topic
	}()

	for {
		select {
		case <-t.ctx.Done():
			return

		case m, ok := <-t.broadcast:
			if !ok {
				return
			}
			payload, err := proto.Marshal(m)
			if err == nil {
				for c := range t.clients {
					c.(interfaces.Client).BroadcastOut(payload)
				}
			}

		case c := <-t.addClient:
			log.Tracef("adding client '%s' to topic '%s'", c.GetUserID(), t.GetName())
			t.clients.Insert(c)

		case c := <-t.removeClient:
			log.Tracef("removing client '%s' to topic '%s'", c.GetUserID(), t.GetName())
			t.clients.Remove(c)

			if t.clients.Size() == 0 {
				return
			}
		}
	}
}
