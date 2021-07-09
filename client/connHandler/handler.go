package connHandler

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	Send      chan *message.Chat
	Uid       string
	conn      *websocket.Conn
	interrupt chan os.Signal
	Done      chan struct{}
}

func NewHandler(uid string) *Handler {
	return &Handler{
		Send:      make(chan *message.Chat),
		Uid:       uid,
		interrupt: make(chan os.Signal),
		Done:      make(chan struct{}),
	}
}

func (h *Handler) Run() {
	h.connectToServer()
	h.handshake()

	go h.writePump()
	go h.readPump()
}

func (h *Handler) connectToServer() {
	endpoint := fmt.Sprintf("127.0.0.1:%d", constants.FrontendClientBasePort)
	u := url.URL{Scheme: "ws", Host: endpoint, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	conn.SetReadLimit(constants.MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(constants.ReadTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(constants.PongTimeout))
		return nil
	})
	h.conn = conn
}

func (h *Handler) handshake() {
	m := &message.Chat{
		Type:   constants.USER,
		UserID: h.Uid,
	}
	err := sendMessage(m, h.conn)
	if err != nil {
		log.Fatal("handshake failed due to ", err)
	}
}

func (h *Handler) writePump() {
	pingTicker := time.NewTicker(constants.PingTimeout)
	defer func() {
		h.conn.Close()
		pingTicker.Stop()
	}()

	for {
		select {
		case <-pingTicker.C:
			h.conn.SetWriteDeadline(time.Now().Add(constants.WriteTimeout))
			if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-h.Done:
			return
		case m := <-h.Send:
			err := sendMessage(m, h.conn)
			if err != nil {
				return
			}
		case <-h.interrupt:
			h.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			select {
			case <-h.Done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (h *Handler) readPump() {
	defer close(h.Done)
	for {
		_, payload, err := h.conn.ReadMessage()
		if err != nil {
			return
		}
		var m message.Chat
		err = proto.Unmarshal(payload, &m)
		if err != nil {
			continue
		}
		log.Printf("[r/%s]> (@%s): %s\n", m.Topic, m.GetUserID(), m.Text)
	}
}

func sendMessage(m *message.Chat, conn *websocket.Conn) error {
	payload, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	conn.SetWriteDeadline(time.Now().Add(constants.WriteTimeout))
	err = conn.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		return err
	}
	return nil
}
