package subscribe

import (
	"log"

	"github.com/afrozalm/minimess/domain/client"
	"github.com/afrozalm/minimess/domain/topic"
)

type service struct {
	tR topic.Repository
}

/* the interface for subscribe should be as simple as below
SubscribeClientToTopic(c *client, t string)
we take a connection and subscribe it to a topic

UnsubscribeClientFromTopic(c *client, t string)
*/

func (s *service) SubscribeClientToTopic(c *client.Client, topicName string) {
	t, err := s.tR.Get(topicName)
	switch err {
	case nil:
		t.AddUserConnectionById(c.Uid, c.Conn)
		c.AddTopicToSubscribedList(topicName)
	case topic.ErrUnknown:
		t = &topic.Topic{Name: topicName}
		t.AddUserConnectionById(c.Uid, c.Conn)
		err := s.tR.Add(t)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(err)
	}
}

func (s *service) UnsubscribeClientFromTopic(c *client.Client, topicName string) {
	t, err := s.tR.Get(topicName)
	switch err {
	case nil:
		t.RemoveUserConnectionById(c.Uid)
		c.RemoveTopicFromSubscribedList(topicName)
	case topic.ErrUnknown:
		log.Fatal("tried to remove a topic which does not exist")
	default:
		log.Fatal(err)
	}
}
