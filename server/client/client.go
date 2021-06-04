package client

/*
	This package contains the client struct which describes the client representation for the server
*/

import (
	"fmt"
	"sync"

	"github.com/afrozalm/minimess/sets"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	uid              string
	Send             chan []byte
	Conn             *websocket.Conn
	SubscribedTopics sets.Set
	mx               *sync.Mutex
}

func NewClient(uid string, conn *websocket.Conn) *Client {
	return &Client{
		uid:              uid,
		Conn:             conn,
		Send:             make(chan []byte, 1024),
		SubscribedTopics: make(sets.Set),
		mx:               &sync.Mutex{},
	}
}

func (c *Client) GetUserID() string {
	return c.uid
}

func (c *Client) AddTopicToSubscribedList(name string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.SubscribedTopics.Insert(name)
	log.Debug(fmt.Sprintf("added '%s' to subscribed topic for client '%s'", name, c.uid))
}

func (c *Client) RemoveTopicFromSubscribedList(name string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.SubscribedTopics.Remove(name)
	log.Debug(fmt.Sprintf("removed '%s' to subscribed topic for client '%s'", name, c.uid))
}

func (c *Client) Close() {
	c.Conn.Close()
	log.Debug(fmt.Sprintf("closing conn '%s'", c.uid))
}
