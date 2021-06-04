package connHandler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/server/client"
	"github.com/afrozalm/minimess/server/supervisor"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWSConn(s *supervisor.Supervisor, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Could not upgrade to ws connection due to", err)
	}
	// now do get the userID and spawn a goroutine
	go parseAndStartPumps(s, conn)
}

func parseAndStartPumps(s *supervisor.Supervisor, conn *websocket.Conn) {
	conn.SetReadLimit(constants.MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(constants.ReadTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(constants.PongTimeout))
		return nil
	})

	// expecting a userID
	m, ok := readMessage(conn)
	if !ok {
		return
	}
	if m.Type != constants.USER {
		log.Info("remote error: first message should be user type")
	}
	log.Debug("got user id ", m.GetUserID())

	c := client.NewClient(m.GetUserID(), conn)
	s.AddClient(c)

	go writePump(s, c)
	go readPump(s, c)
}

func writePump(s *supervisor.Supervisor, c *client.Client) {
	pingTicker := time.NewTicker(constants.PingTimeout)
	defer func() {
		pingTicker.Stop()
		log.Trace(fmt.Sprintf("closing readPump for user '%s'", c.GetUserID()))
		c.Close()
		s.OnClientClose(c)
	}()
	for {
		select {
		case <-pingTicker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteTimeout))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			log.Trace(fmt.Sprintf("ping sent to client '%s'", c.GetUserID()))
		case m, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(constants.WriteTimeout))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Trace(fmt.Sprintf("closing connection for client '%s'", c.GetUserID()))
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
			log.Trace(fmt.Sprintf("message sent to '%s'", c.GetUserID()))
		}
	}
}

func readPump(s *supervisor.Supervisor, c *client.Client) {
	defer func() {
		log.Trace(fmt.Sprintf("closing readPump for user '%s'", c.GetUserID()))
		c.Close()
	}()

	for {
		m, ok := readMessage(c.Conn)
		if !ok {
			return
		}
		switch m.Type {
		case constants.SUBSCRIBE:
			log.Trace(fmt.Sprintf("got subscribe request from '%s' to '%s'", c.GetUserID(), m.GetTopic()))
			s.SubscribeClientToTopic(c, m.Topic)
		case constants.UNSUBSCRIBE:
			log.Trace(fmt.Sprintf("got unsubscribe request from '%s' to '%s'", c.GetUserID(), m.GetTopic()))
			s.UnsubscribeClientFromTopic(c, m.Topic)
		case constants.CHAT:
			log.Trace(fmt.Sprintf("going to broadcast '%s' to '%s'", m.Text, m.GetTopic()))
			m.TimeUUID = gocql.TimeUUID().Bytes()
			s.BroadcastMessageToKafka(m)
			log.Trace(fmt.Sprintf("broadcasted '%s' to '%s'", m.GetText(), m.GetTopic()))
		default:
			log.Warn(fmt.Sprintf("not handled message for type: '%s', %v", m.GetType(), m))
		}
	}
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
