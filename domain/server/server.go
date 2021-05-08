package server

import (
	"sync"

	"github.com/afrozalm/minimess/domain/client"
	"github.com/afrozalm/minimess/domain/topic"
	"github.com/afrozalm/minimess/message"
)

/*
	the server is supposed to listen to ws connections on localhost:9091
	for each new connection, you read what the client id is and create a client
	start a goroutine to handle the client.
	The server also stores topics in memory thus we need to create a storage object attached
	server
*/

// type clientSet map[*client.Client]struct{}

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
	server.topicMx.Lock()
	defer server.topicMx.Unlock()
	server.topics[topic.Name] = topic
}

func (server *Server) SubscribeClientToTopic(client *client.Client, name string) {
	var t *topic.Topic
	if !server.existsTopic(name) {
		t = topic.NewTopic(name)
		t.Run()
		server.addTopic(t)
	} else {
		t = server.getTopic(name)
	}
	t.AddClient <- client
	client.AddTopicToSubscribedList(name)
}

func (server *Server) UnsubscribeClientFromTopic(client *client.Client, name string) {
	if !server.existsTopic(name) {
		return
	}
	topic := server.getTopic(name)
	topic.RemoveClient <- client
	client.RemoveTopicFromSubscribedList(name)
}

func (server *Server) BroadcastMessageToTopic(message *message.Message) {
	topic := server.getTopic(message.Topic)
	if topic == nil {
		return
	}
	topic.Broadcast <- message
}
