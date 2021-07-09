package broadcaster

import (
	"context"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/server/fanout/interfaces"
	frontend "github.com/afrozalm/minimess/server/frontend/frontend_proto"
	"github.com/afrozalm/minimess/server/gokart/consumer"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Broadcaster struct {
	consumer            *consumer.AtLeastOnceConsumer
	storage             interfaces.Storage
	receivedMessages    chan kafka.Message
	broadcastedMessages chan kafka.Message
	frontend_stub_map   map[string]frontend.FrontEndClient
	ctx                 context.Context
}

func NewBroadcaster(ctx context.Context, storage interfaces.Storage) *Broadcaster {
	consumer := consumer.NewAtLeastOnceConsumer(ctx, constants.BroadcasterCG)
	receivedMessages := make(chan kafka.Message, 10)
	broadcastedMessages := make(chan kafka.Message, 20)
	go consumer.Run(receivedMessages, broadcastedMessages)

	return &Broadcaster{
		consumer:            consumer,
		storage:             storage,
		receivedMessages:    receivedMessages,
		broadcastedMessages: broadcastedMessages,
		frontend_stub_map:   make(map[string]frontend.FrontEndClient),
		ctx:                 ctx,
	}
}

func (broadcaster *Broadcaster) Run() {
	defer func() {
		log.Info("closing broadcaster's broadcastedMessages channel updates")
		close(broadcaster.broadcastedMessages)
	}()
	for {
		select {
		case m, ok := <-broadcaster.receivedMessages:
			if !ok {
				return
			}
			topic := string(m.Key)
			chatmessage := &message.Chat{}
			err := proto.Unmarshal(m.Value, chatmessage)
			if err != nil {
				log.Warnf("failed to unmarshall chat message", m.Offset, m.Topic)
			}
			hosts := broadcaster.storage.GetHostsForTopic(topic)
			for _, host := range hosts {
				if stub, ok := broadcaster.frontend_stub_map[host]; ok {
					stub.BroadcastOut(broadcaster.ctx, chatmessage)
				} else {
					conn, err := grpc.Dial(host, grpc.WithInsecure())
					if err != nil {
						log.Warnf("could not create a grpc connection for host '%s'", host)
						continue
					}
					defer conn.Close()
					stub = frontend.NewFrontEndClient(conn)
					broadcaster.frontend_stub_map[host] = stub
					stub.BroadcastOut(broadcaster.ctx, chatmessage)
					log.Debugf("broadcasted out message '%s'", chatmessage.Text)
				}
			}
			broadcaster.broadcastedMessages <- m
		case <-broadcaster.ctx.Done():
			return
		}
	}
}
