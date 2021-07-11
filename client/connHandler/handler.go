package connHandler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	Send chan *message.Chat
	Uid  string
	conn *websocket.Conn
}

func NewHandler(uid string) *Handler {
	return &Handler{
		Send: make(chan *message.Chat),
		Uid:  uid,
	}
}

func (h *Handler) Run(ctx context.Context) {
	for {
		h.connectToServer()
		h.handshake()

		go h.readPump(ctx)
		h.writePump(ctx)

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (h *Handler) connectToServer() {
	endpoint := fmt.Sprintf("127.0.0.1:%d", constants.FrontendClientBasePort)
	u := url.URL{Scheme: "ws", Host: endpoint, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("err connection to server:", err)
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

func (h *Handler) writePump(ctx context.Context) {
	pingTicker := time.NewTicker(constants.PingTimeout)
	defer func() {
		h.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
		case m := <-h.Send:
			err := sendMessage(m, h.conn)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Handler) readPump(ctx context.Context) {
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

		select {
		case <-ctx.Done():
			return
		default:
		}
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
