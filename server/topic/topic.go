package topic

import (
	"fmt"

	"github.com/afrozalm/minimess/server/client"
	"github.com/afrozalm/minimess/sets"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type Topic struct {
	clients      sets.Set
	Name         string
	broadcast    chan kafka.Message
	addClient    chan *client.Client
	removeClient chan *client.Client
	done         chan struct{}
}

func NewTopic(name string) *Topic {
	// TODO make an rpc to connection service and subscribe to topic
	return &Topic{
		clients:      make(sets.Set),
		Name:         name,
		broadcast:    make(chan kafka.Message, 10),
		addClient:    make(chan *client.Client, 10),
		removeClient: make(chan *client.Client, 10),
		done:         make(chan struct{}, 1),
	}
}

func (t *Topic) Close() {
	log.Trace("closing topic goroutine for topic ", t.GetName())
	t.done <- struct{}{}
	close(t.done)
	close(t.addClient)
	close(t.broadcast)
	close(t.removeClient)
	log.Trace("closed topic ", t.GetName())
}

func (t *Topic) RemoveClient(c *client.Client) {
	t.removeClient <- c
}

func (t *Topic) AddClient(c *client.Client) {
	t.addClient <- c
}

func (t *Topic) Broadcast(m kafka.Message) {
	t.broadcast <- m
}

func (t *Topic) GetName() string {
	return t.Name
}

func (t *Topic) Run(broadcastedMessages chan kafka.Message) {
	defer func() {
		// garbage collect
		t.clients = nil
	}()
	for {
		select {
		case <-t.done:
			return

		case m, ok := <-t.broadcast:
			if !ok {
				return
			}
			for c := range t.clients {
				c.(*client.Client).Send <- m.Value
			}
			broadcastedMessages <- m

		case c := <-t.addClient:
			log.Trace(fmt.Sprintf("adding client '%s' to topic '%s'", c.GetUserID(), t.Name))
			t.clients.Insert(c)

		case c := <-t.removeClient:
			log.Trace(fmt.Sprintf("removing client '%s' to topic '%s'", c.GetUserID(), t.Name))
			t.clients.Remove(c)
		}
	}
}
