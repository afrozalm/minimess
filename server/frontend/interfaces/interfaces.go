package interfaces

import (
	"context"

	"github.com/afrozalm/minimess/message"
)

type Client interface {
	GetUserID() string

	AddTopic(string)

	RemoveTopic(string)

	BroadcastOut([]byte)
}

type Supervisor interface {
	GetCtx() context.Context

	SubscribeClientToTopic(Client, string)

	UnsubscribeClientFromTopic(Client, string)

	BroadcastIn(*message.Chat)
}

type Topic interface {
	GetName() string

	AddClient(Client)

	RemoveClient(Client)

	BroadcastOut(*message.Chat)
}
