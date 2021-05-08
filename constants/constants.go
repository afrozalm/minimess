package constants

import (
	"time"
)

const (
	// timeouts in seconds
	PongTimeout  = 60 * time.Second
	ReadTimeout  = PongTimeout
	PingTimeout  = (PongTimeout * 4) / 5
	WriteTimeout = 10 * time.Second

	// buffer limits
	MaxMessageSize = 512

	// message constants
	// Types
	USER        = "user"
	SUBSCRIBE   = "sub"
	UNSUBSCRIBE = "unsub"
	CHAT        = "send"
)
