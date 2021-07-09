package supervisor

import (
	"context"

	"github.com/afrozalm/minimess/message"
	fanout "github.com/afrozalm/minimess/server/fanout/fanout_proto"
	"github.com/afrozalm/minimess/server/frontend/interfaces"
	"github.com/afrozalm/minimess/server/frontend/topic"
	log "github.com/sirupsen/logrus"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Supervisor) GetCtx() context.Context {
	return s.ctx
}

func (s *Supervisor) SubscribeClientToTopic(c interfaces.Client, topicName string) {
	var t topic.Topic
	log.Tracef("subscribing '%s' to '%s'", c.GetUserID(), topicName)
	if !s.existsTopic(topicName) {
		t = *topic.NewTopic(topicName, s.GetCtx())
		go t.Run(s.closedTopics)
		s.addTopic(t)
	} else {
		t = s.getTopic(topicName).(topic.Topic)
	}
	t.AddClient(c)
	c.AddTopic(topicName)
	log.Tracef("subscribed '%s' to '%s'", c.GetUserID(), topicName)
}

func (s *Supervisor) UnsubscribeClientFromTopic(client interfaces.Client, name string) {
	if !s.existsTopic(name) {
		log.Tracef("topic does not exist to unsubscribe client '%s' from '%s'", client.GetUserID(), name)
		return
	}
	topic := s.getTopic(name)
	topic.RemoveClient(client)
	client.RemoveTopic(name)
	log.Tracef("unsubscribed '%s' from '%s'", client.GetUserID(), name)
}

func (s *Supervisor) BroadcastOut(ctx context.Context, message *message.Chat) (*emptypb.Empty, error) {
	if ctx.Err() == nil {
		s.receivedMessages <- message
		return new(emptypb.Empty), nil
	}
	if ctx.Err() == context.Canceled {
		return new(emptypb.Empty), status.Errorf(codes.Canceled, "context is cancelled")
	}
	if ctx.Err() == context.DeadlineExceeded {
		return new(emptypb.Empty), status.Errorf(codes.DeadlineExceeded, "context deadline exceeded")
	}
	return new(emptypb.Empty), ctx.Err()
}

func (s *Supervisor) BroadcastIn(message *message.Chat) {
	err := s.messageProducer.Produce(message)
	if err != nil {
		log.Infof("failed to broadcast message '%s'", message.GetText())
	}
}

func (s *Supervisor) refreshSubscriptions() {
	// TODO: call rpc to fanout and refresh ttl on all topics
	for topic := range s.topics {
		s.fanout_stub.Subscribe(s.ctx, &fanout.HostTopic{Host: s.addr, Topic: topic})
	}
	log.Info("refreshed all topic subs")
}

func (s *Supervisor) existsTopic(name string) bool {
	s.topicMx.RLock()
	defer s.topicMx.RUnlock()
	_, ok := s.topics[name]
	return ok
}

func (s *Supervisor) getTopic(name string) interfaces.Topic {
	s.topicMx.RLock()
	defer s.topicMx.RUnlock()
	if topic, ok := s.topics[name]; ok {
		return topic
	}
	return nil
}

func (s *Supervisor) addTopic(topic interfaces.Topic) {
	log.Tracef("adding new topic: '%s' to supervisor", topic.GetName())
	s.topicMx.Lock()
	defer s.topicMx.Unlock()
	s.topics[topic.GetName()] = topic
	s.subscribeTopicToFanout(topic.GetName())
}

func (s *Supervisor) removeTopic(topic interfaces.Topic) {
	log.Tracef("closing topic '%s'", topic.GetName())
	s.topicMx.Lock()
	defer s.topicMx.Unlock()
	delete(s.topics, topic.GetName())
	s.unsubscribeTopicFromFanout(topic.GetName())
}

func (s *Supervisor) subscribeTopicToFanout(topic string) {
	s.fanout_stub.Subscribe(s.ctx, &fanout.HostTopic{Host: s.addr, Topic: topic})
}

func (s *Supervisor) unsubscribeTopicFromFanout(topic string) {
	s.fanout_stub.Unsubscribe(s.ctx, &fanout.HostTopic{Host: s.addr, Topic: topic})
}

func (s *Supervisor) unsubscribeFromFanout() {
	s.fanout_stub.UnsubscribeAll(context.Background(), &fanout.Host{Host: s.addr})
}
