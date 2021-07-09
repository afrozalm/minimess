package connHandler

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/server/frontend/client"
	"github.com/afrozalm/minimess/server/frontend/supervisor"
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

	go c.RunWritePump(s)
	go c.RunReadPump(s)
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
