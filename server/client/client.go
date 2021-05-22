package client

import (
	"sync"

	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/sets"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	Uid              string
	Send             chan *message.Message
	Conn             *websocket.Conn
	SubscribedTopics sets.Set
	mx               *sync.Mutex
}

func NewClient(uid string, conn *websocket.Conn) *Client {
	return &Client{
		Uid:              uid,
		Conn:             conn,
		Send:             make(chan *message.Message, 10),
		SubscribedTopics: make(sets.Set),
		mx:               &sync.Mutex{},
	}
}

func (c *Client) AddTopicToSubscribedList(name string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.SubscribedTopics.Insert(name)
	log.Debug("added '%s' to subscribed topic for client '%s'", name, c.Uid)
}

func (c *Client) RemoveTopicFromSubscribedList(name string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.SubscribedTopics.Remove(name)
	log.Debug("removed '%s' to subscribed topic for client '%s'", name, c.Uid)
}

func (c *Client) Close() {
	c.Conn.Close()
	log.Debug("closing conn '%s'", c.Uid)
}
