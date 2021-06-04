package supervisor

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/afrozalm/minimess/message"
	"github.com/afrozalm/minimess/server/client"
	"github.com/afrozalm/minimess/server/gokart/consumer"
	"github.com/afrozalm/minimess/server/gokart/producer"
	"github.com/afrozalm/minimess/server/topic"
	"github.com/afrozalm/minimess/sets"
)

type Supervisor struct {
	topicMx             *sync.RWMutex
	topics              map[string]*topic.Topic
	clients             sets.Set
	clientMx            *sync.RWMutex
	messageProducer     *producer.AtLeastOnceProducer
	messageConsumer     *consumer.AtLeastOnceConsumer
	receivedMessages    chan kafka.Message
	broadcastedMessages chan kafka.Message
	done                chan struct{}
}

func NewSupervisor() *Supervisor {
	messageConsumer := consumer.NewAtLeastOnceConsumer(fmt.Sprintf("cg-test-%d", rand.Int()))
	receivedMessages := make(chan kafka.Message, 10)
	broadcastedMessages := make(chan kafka.Message, 10)
	go messageConsumer.Run(receivedMessages, broadcastedMessages)
	log.Trace("create a new supervisor")
	return &Supervisor{
		topicMx:         &sync.RWMutex{},
		topics:          make(map[string]*topic.Topic),
		clients:         make(sets.Set),
		clientMx:        &sync.RWMutex{},
		messageProducer: producer.NewAtLeaseOnceProducer(-1, 3),
		// for now, all servers belong to different consumer groups
		// I know, not the best thing to will. will be fixed with a fanout service
		messageConsumer:     messageConsumer,
		receivedMessages:    receivedMessages,
		broadcastedMessages: broadcastedMessages,
		done:                make(chan struct{}, 1),
	}
}

func (supervisor *Supervisor) Run() {
	log.Trace("running supervisor")
	for {
		select {
		case <-supervisor.done:
			return
		case m, ok := <-supervisor.receivedMessages:
			if !ok {
				log.Warn("receivedMessage chan is closed. exiting supervisor process run")
				return
			}
			topicName := string(m.Key)
			topic := supervisor.getTopic(topicName)
			if topic == nil {
				log.Trace(fmt.Sprintf("topic '%s' does not exist to broadcast message", topicName))
				continue
			}
			topic.Broadcast(m)
		}
	}
}

func (supervisor *Supervisor) Close() {
	log.Trace("closing supervisor")
	supervisor.messageConsumer.Close()
	supervisor.messageProducer.Close()
	// close all clients connections
	log.Trace("going to close clients")
	for c := range supervisor.clients {
		c.(*client.Client).Close()
	}
	supervisor.clients = nil
	log.Trace("clients closed. going to close all topics")
	// close all topics
	for _, t := range supervisor.topics {
		t.Close()
	}
	supervisor.topics = nil
	// close all channels
	close(supervisor.broadcastedMessages)
	log.Trace("closed all topics and broadcastedMessages channel")
	supervisor.done <- struct{}{}
	close(supervisor.done)
	log.Trace("supervisor close done")
}

func (supervisor *Supervisor) OnTopicClose(t *topic.Topic) {
	name := t.Name
	log.Trace(fmt.Sprintf("closing topic '%s'", t.Name))
	t.Close()
	supervisor.topicMx.Lock()
	defer supervisor.topicMx.Unlock()
	delete(supervisor.topics, name)
}

func (supervisor *Supervisor) OnClientClose(c *client.Client) {
	log.Trace(fmt.Sprintf("removing client '%s' from all subscribed topics", c.GetUserID()))
	for name := range c.SubscribedTopics {
		t := supervisor.getTopic(name.(string))
		if t != nil {
			t.RemoveClient(c)
		}
	}
}

func (supervisor *Supervisor) existsTopic(name string) bool {
	supervisor.topicMx.RLock()
	defer supervisor.topicMx.RUnlock()
	_, ok := supervisor.topics[name]
	return ok
}

func (supervisor *Supervisor) getTopic(name string) *topic.Topic {
	supervisor.topicMx.RLock()
	defer supervisor.topicMx.RUnlock()
	if topic, ok := supervisor.topics[name]; ok {
		return topic
	}
	return nil
}

func (supervisor *Supervisor) AddClient(client *client.Client) {
	supervisor.clientMx.Lock()
	defer supervisor.clientMx.Unlock()
	supervisor.clients.Insert(client)
}

func (supervisor *Supervisor) addTopic(topic *topic.Topic) {
	log.Trace(fmt.Sprintf("adding new topic: '%s' to supervisor", topic.Name))
	supervisor.topicMx.Lock()
	defer supervisor.topicMx.Unlock()
	supervisor.topics[topic.Name] = topic
}

func (supervisor *Supervisor) SubscribeClientToTopic(client *client.Client, name string) {
	var t *topic.Topic
	log.Trace(fmt.Sprintf("subscribing '%s' to '%s'", client.GetUserID(), name))
	if !supervisor.existsTopic(name) {
		t = topic.NewTopic(name)
		go t.Run(supervisor.broadcastedMessages)
		supervisor.addTopic(t)
	} else {
		t = supervisor.getTopic(name)
	}
	t.AddClient(client)
	client.AddTopicToSubscribedList(name)
	log.Trace(fmt.Sprintf("subscribed '%s' to '%s'", client.GetUserID(), name))
}

func (supervisor *Supervisor) UnsubscribeClientFromTopic(client *client.Client, name string) {
	if !supervisor.existsTopic(name) {
		log.Trace(fmt.Sprintf("topic does not exist to unsubscribe client '%s' from '%s'", client.GetUserID(), name))
		return
	}
	topic := supervisor.getTopic(name)
	topic.RemoveClient(client)
	client.RemoveTopicFromSubscribedList(name)
	log.Trace(fmt.Sprintf("unsubscribed '%s' from '%s'", client.GetUserID(), name))
}

func (supervisor *Supervisor) BroadcastMessageToKafka(message *message.Chat) {
	supervisor.messageProducer.Produce(message)
}
