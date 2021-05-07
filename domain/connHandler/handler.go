package connHandler

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/afrozalm/minimess/domain/client"
	"github.com/afrozalm/minimess/domain/constants"
	"github.com/afrozalm/minimess/domain/message"
	"github.com/afrozalm/minimess/domain/server"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWSConn(s *server.Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Could not upgrade to ws connection due to", err)
	}
	// now do get the userID and spawn a goroutine
	go parseAndStartPumps(s, conn)
}

func parseAndStartPumps(s *server.Server, conn *websocket.Conn) {
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
		log.Println("remote error: first message should be user type")
	}
	log.Println("got user id", m.Uid)

	c := client.NewClient(m.Uid, conn)

	go writePump(s, c)
	go readPump(s, c)
}

func writePump(s *server.Server, c *client.Client) {
	pingTicker := time.NewTicker(constants.PingTimeout)
	defer func() {
		pingTicker.Stop()
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
		case m, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writeMessage(m, w)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				writeMessage(<-c.Send, w)
			}
		}
	}
}

func readPump(s *server.Server, c *client.Client) {
	defer func() {
		c.Close()
	}()

	for {
		m, ok := readMessage(c.Conn)
		if !ok {
			return
		}
		switch m.Type {
		case constants.SUBSCRIBE:
			s.SubscribeClientToTopic(c, m.Topic)
		case constants.UNSUBSCRIBE:
			s.UnsubscribeClientFromTopic(c, m.Topic)
		case constants.CHAT:
			s.BroadcastMessageToTopic(m)
		default:
			log.Printf("not handled message for type: %s, %v", m.Type, m)
		}
	}
}

func readMessage(conn *websocket.Conn) (*message.Message, bool) {
	_, payload, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("closing connection with err:", err)
			return nil, false
		}
	}

	m, err := message.DecodeMessage(payload)
	if err != nil {
		log.Println("remote error: closing due to bad handshake", err)
		return m, false
	}
	return m, true
}

func writeMessage(m *message.Message, w io.WriteCloser) {
	payload, err := m.EncodeMessage()
	if err != nil {
		log.Printf("bad message not forwarding %v", *m)
	}
	w.Write(payload)

}
