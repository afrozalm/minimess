package topic

import (
	"github.com/afrozalm/minimess/domain/client"
	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/sets"
	log "github.com/sirupsen/logrus"
)

type Topic struct {
	Name         string
	Broadcast    chan *message.Message
	clients      sets.Set
	Close        chan bool
	AddClient    chan *client.Client
	RemoveClient chan *client.Client
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name:         name,
		Broadcast:    make(chan *message.Message, 10),
		clients:      make(sets.Set),
		Close:        make(chan bool),
		AddClient:    make(chan *client.Client, 10),
		RemoveClient: make(chan *client.Client, 10),
	}
}

func (t *Topic) Run() {
	for {
		select {
		case m := <-t.Broadcast:
			for c := range t.clients {
				c.(*client.Client).Send <- m
			}
		case c := <-t.AddClient:
			log.Trace("adding client '%s' to topic '%s'", c.Uid, t.Name)
			t.clients.Insert(c)
		case c := <-t.RemoveClient:
			log.Trace("removing client '%s' to topic '%s'", c.Uid, t.Name)
			t.clients.Remove(c)
		case <-t.Close:
			log.Trace("closing topic '%s'", t.Name)
			return
		}
	}
}
