package client

import (
	"io"
	"time"

	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/server/frontend/interfaces"
	"github.com/gorilla/websocket"
)

func (c *Client) RunWritePump(s interfaces.Supervisor) {
	pingTicker := time.NewTicker(constants.PingTimeout)
	defer func() {
		pingTicker.Stop()
		log.Tracef("closing writePump for user '%s'", c.GetUserID())
		c.close(s)
	}()

	for {
		select {
		case <-pingTicker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteTimeout))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			log.Tracef("ping sent to client '%s'", c.GetUserID())
		case m, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteTimeout))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Tracef("closing connection for client '%s'", c.GetUserID())
				return
			}
			w, err := c.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			writeMessage(m, w)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				writeMessage(<-c.Send, w)
			}
			if err := w.Close(); err != nil {
				return
			}
			log.Tracef("message sent to '%s'", c.GetUserID())
		case <-s.GetCtx().Done():
			log.Trace("parent context closed, closing connection")
			return
		}
	}
}

func (c *Client) RunReadPump(s interfaces.Supervisor) {
	defer func() {
		log.Tracef("closing readPump for user '%s'", c.GetUserID())
	}()

	for {
		if s.GetCtx().Err() != nil {
			return
		}

		m, ok := readMessage(c.Conn)
		if !ok {
			return
		}
		switch m.Type {
		case constants.SUBSCRIBE:
			log.Tracef("got subscribe request from '%s' to '%s'", c.GetUserID(), m.GetTopic())
			s.SubscribeClientToTopic(c, m.GetTopic())
		case constants.UNSUBSCRIBE:
			log.Tracef("got unsubscribe request from '%s' to '%s'", c.GetUserID(), m.GetTopic())
			s.UnsubscribeClientFromTopic(c, m.GetTopic())
		case constants.CHAT:
			log.Tracef("going to broadcast in '%s' to '%s'", m.Text, m.GetTopic())
			m.TimeUUID = gocql.TimeUUID().Bytes()
			s.BroadcastIn(m)
			log.Tracef("broadcasted in '%s' to '%s'", m.GetText(), m.GetTopic())
		default:
			log.Warnf("not handled message for type: '%s', %v", m.GetType(), m)
		}
	}

}

func (c *Client) BroadcastOut(payload []byte) {
	c.Send <- payload
}

func (c *Client) GetUserID() string {
	return c.uid
}

func (c *Client) AddTopic(name string) {
	c.subscribedTopics.Insert(name)
	log.Debugf("added '%s' to subscribed topic for client '%s'", name, c.uid)
}

func (c *Client) RemoveTopic(name string) {
	c.subscribedTopics.Remove(name)
	log.Debugf("removed '%s' to subscribed topic for client '%s'", name, c.uid)
}

func (c *Client) close(s interfaces.Supervisor) {
	c.Conn.Close()
	for topic := range c.subscribedTopics {
		s.UnsubscribeClientFromTopic(c, topic.(string))
	}
	log.Debugf("closing conn '%s'", c.uid)
}

func readMessage(conn *websocket.Conn) (*message.Chat, bool) {
	_, payload, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Info("closing connection with err:", err)
			return nil, false
		}
		log.Info("read error from connection")
		return nil, false
	}
	var m message.Chat
	err = proto.Unmarshal(payload, &m)
	if err != nil {
		log.Warn("remote error: closing due to bad payload", err)
		return nil, false
	}
	return &m, true
}

func writeMessage(payload []byte, w io.WriteCloser) {
	w.Write(payload)
}
