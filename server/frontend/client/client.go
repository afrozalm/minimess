package client

import (
	"github.com/afrozalm/minimess/sets"
	"github.com/gorilla/websocket"
)

type Client struct {
	uid              string
	Send             chan []byte
	Conn             *websocket.Conn
	subscribedTopics sets.Set
}

func NewClient(uid string, conn *websocket.Conn) *Client {
	return &Client{
		uid:              uid,
		Conn:             conn,
		Send:             make(chan []byte, 1024),
		subscribedTopics: make(sets.Set),
	}
}
