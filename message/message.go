package message

import (
	"bytes"
	"encoding/json"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Message struct {
	Type  string
	Uid   string
	Topic string
	Text  string
}

func NewMessage() *Message {
	return &Message{}
}

func DecodeMessage(payload []byte) (*Message, error) {
	payload = bytes.TrimSpace(bytes.Replace(payload, newline, space, -1))
	m := &Message{}
	err := json.Unmarshal(payload, m)
	return m, err
}

func (m *Message) EncodeMessage() ([]byte, error) {
	return json.Marshal(m)
}
