package connHandler

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Send      chan *message.Message
	Uid       string
	conn      *websocket.Conn
	interrupt chan os.Signal
	done      chan struct{}
}

func NewHandler(uid string) *Handler {
	return &Handler{
		Send:      make(chan *message.Message),
		Uid:       uid,
		interrupt: make(chan os.Signal),
		done:      make(chan struct{}),
	}
}

func (h *Handler) Run() {
	h.connectToServer("127.0.0.1", "8080")
	h.handshake()

	go h.writePump()
	go h.readPump()
}

func (h *Handler) connectToServer(ip string, port string) {
	endpoint := ip + ":" + port
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
	m := new(message.Message)
	m.Type = constants.USER
	m.Uid = h.Uid
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
		case <-h.done:
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
			case <-h.done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (h *Handler) readPump() {
	defer close(h.done)
	for {
		_, payload, err := h.conn.ReadMessage()
		if err != nil {
			return
		}
		m, err := message.DecodeMessage(payload)
		if err != nil {
			continue
		}
		log.Printf("[r/%s]> (@%s): %s\n", m.Topic, m.Uid, m.Text)
	}
}

func sendMessage(m *message.Message, conn *websocket.Conn) error {
	payload, err := m.EncodeMessage()
	if err != nil {
		return err
	}
	conn.SetWriteDeadline(time.Now().Add(constants.WriteTimeout))
	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return err
	}
	return nil
}
