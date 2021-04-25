package client

import "github.com/gorilla/websocket"

type stringSet map[string]struct{}

type Client struct {
	Uid              string
	Conn             *websocket.Conn
	SubscribedTopics stringSet
}

func (s stringSet) insert(key string) {
	s[key] = struct{}{}
}

func (s stringSet) remove(key string) {
	if _, ok := s[key]; ok {
		delete(s, key)
	}
}

func (s stringSet) exists(key string) bool {
	if _, ok := s[key]; ok {
		return true
	}
	return false
}

func GetClient(uid string, conn *websocket.Conn) *Client {
	return &Client{Uid: uid, Conn: conn}
}

func (c *Client) AddTopicToSubscribedList(topic string) {
	c.SubscribedTopics.insert(topic)
}

func (c *Client) RemoveTopicFromSubscribedList(topic string) {
	c.SubscribedTopics.remove(topic)
}
