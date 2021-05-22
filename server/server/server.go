package server

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/server/client"
	"github.com/afrozalm/minimess/server/topic"
)

type Server struct {
	topics  map[string]*topic.Topic
	topicMx *sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		topics:  make(map[string]*topic.Topic),
		topicMx: &sync.RWMutex{},
	}
}

func (server *Server) OnClientClose(c *client.Client) {
	log.Trace("removing client '%s' from all subscribed topics", c.Uid)
	for name := range c.SubscribedTopics {
		t := server.getTopic(name.(string))
		if t != nil {
			t.RemoveClient <- c
		}
	}
}

func (server *Server) existsTopic(name string) bool {
	server.topicMx.RLock()
	defer server.topicMx.RUnlock()
	_, ok := server.topics[name]
	return ok
}

func (server *Server) getTopic(name string) *topic.Topic {
	server.topicMx.RLock()
	defer server.topicMx.RUnlock()
	if topic, ok := server.topics[name]; ok {
		return topic
	}
	return nil
}

func (server *Server) addTopic(topic *topic.Topic) {
	log.Trace("adding new topic: '%s' to server", topic.Name)
	server.topicMx.Lock()
	defer server.topicMx.Unlock()
	server.topics[topic.Name] = topic
}

func (server *Server) SubscribeClientToTopic(client *client.Client, name string) {
	var t *topic.Topic
	log.Trace("subscribing '%s' to '%s'", client.Uid, name)
	if !server.existsTopic(name) {
		t = topic.NewTopic(name)
		go t.Run()
		server.addTopic(t)
	} else {
		t = server.getTopic(name)
	}
	t.AddClient <- client
	client.AddTopicToSubscribedList(name)
	log.Trace("subscribed '%s' to '%s'", client.Uid, name)
}

func (server *Server) UnsubscribeClientFromTopic(client *client.Client, name string) {
	if !server.existsTopic(name) {
		log.Trace("topic does not exist to unsubscribe client '%s' from '%s'", client.Uid, name)
		return
	}
	topic := server.getTopic(name)
	topic.RemoveClient <- client
	client.RemoveTopicFromSubscribedList(name)
	log.Trace("unsubscribed '%s' from '%s'", client.Uid, name)
}

func (server *Server) BroadcastMessageToTopic(message *message.Message) {
	topic := server.getTopic(message.Topic)
	if topic == nil {
		log.Trace("topic '%s' does not exist to broadcast message", message.Topic)
		return
	}
	topic.Broadcast <- message
}
