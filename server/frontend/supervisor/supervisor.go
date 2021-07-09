package supervisor

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	fanout "github.com/afrozalm/minimess/server/fanout/fanout_proto"
	frontend "github.com/afrozalm/minimess/server/frontend/frontend_proto"
	"github.com/afrozalm/minimess/server/frontend/interfaces"
	"github.com/afrozalm/minimess/server/gokart/producer"
)

type Supervisor struct {
	topicMx          *sync.RWMutex
	topics           map[string]interfaces.Topic
	messageProducer  *producer.AtLeastOnceProducer
	receivedMessages chan *message.Chat
	closedTopics     chan interfaces.Topic
	ctx              context.Context
	fanout_stub      fanout.FanoutClient
	fanout_conn      *grpc.ClientConn
	addr             string
	frontend.UnimplementedFrontEndServer
}

func NewSupervisor(ctx context.Context, addr string) *Supervisor {
	fanoutAddr := fmt.Sprintf("127.0.0.1:%d", constants.FanoutBasePort)
	conn, err := grpc.Dial(fanoutAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to fanout service '%s' due to '%v'", fanoutAddr, err)
	}
	stub := fanout.NewFanoutClient(conn)
	log.Trace("created a new supervisor")
	return &Supervisor{
		topicMx:          &sync.RWMutex{},
		topics:           make(map[string]interfaces.Topic),
		messageProducer:  producer.NewAtLeaseOnceProducer(ctx, -1, 3),
		receivedMessages: make(chan *message.Chat, 10),
		closedTopics:     make(chan interfaces.Topic, 10),
		ctx:              ctx,
		fanout_stub:      stub,
		fanout_conn:      conn,
		addr:             addr,
	}
}

func (supervisor *Supervisor) Run() {
	log.Trace("running supervisor")
	ticker := time.NewTicker(time.Minute * 5)
	defer func() {
		ticker.Stop()
		log.Trace("closing supervisor")
		supervisor.messageProducer.Close()
		supervisor.topics = nil
		// close all channels
		log.Trace("closed broadcastedMessages channel")
		supervisor.unsubscribeFromFanout()
		supervisor.fanout_conn.Close()
	}()

	for {
		select {
		case <-supervisor.ctx.Done():
			return
		case m, ok := <-supervisor.receivedMessages:
			if !ok {
				log.Warn("receivedMessage chan is closed. exiting supervisor process run")
				return
			}
			topicName := string(m.GetTopic())
			topic := supervisor.getTopic(topicName)
			if topic == nil {
				log.Tracef("topic '%s' does not exist to broadcast message", topicName)
				continue
			}
			topic.BroadcastOut(m)
		case <-ticker.C:
			supervisor.refreshSubscriptions()
		case t := <-supervisor.closedTopics:
			supervisor.removeTopic(t)
		}
	}
}
