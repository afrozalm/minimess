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
	ATTACH      = "attach"
	DETACH      = "//"
	NORMAL      = "normal"
	QUIT        = "quit"

	// kafka
	BootstrapServers = "localhost:9091,localhost:9092,localhost:9093"
	MessageTopic     = "message-log"

	// Broadcaster
	BroadcasterCG = "BROADCASTER-CG"
)

var (
	// kafka
	Brokers = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

	// cassandra
	CassandraHosts = []string{"localhost:9042"}

	// fanout
	FanoutBasePort = 11000

	// frontend
	FrontendClientBasePort = 12000
	FrontendGRPCBasePort   = 13000
)
