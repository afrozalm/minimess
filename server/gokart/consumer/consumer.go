package consumer

import (
	"context"
	"time"

	"github.com/afrozalm/minimess/constants"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type AtLeastOnceConsumer struct {
	Reader *kafka.Reader
	done   chan struct{}
}

func NewAtLeastOnceConsumer(consumerGroup string) *AtLeastOnceConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        constants.Brokers,
		GroupID:        consumerGroup,
		Topic:          constants.MessageTopic,
		IsolationLevel: kafka.ReadCommitted,
		CommitInterval: time.Second,
	})
	return &AtLeastOnceConsumer{
		Reader: reader,
		done:   make(chan struct{}, 1),
	}
}

func (consumer *AtLeastOnceConsumer) consume() (kafka.Message, error) {
	return consumer.Reader.FetchMessage(context.Background())
}

func (consumer *AtLeastOnceConsumer) commitMessages(msgs ...kafka.Message) error {
	return consumer.Reader.CommitMessages(context.Background(), msgs...)
}

func (consumer *AtLeastOnceConsumer) Close() {
	consumer.Reader.Close()
	consumer.done <- struct{}{}
}

func (consumer *AtLeastOnceConsumer) Run(receivedMessages, broadcastedMessages chan kafka.Message) {
	go consumer.committer(broadcastedMessages)
	defer func() {
		close(receivedMessages)
	}()
	for {
		m, e := consumer.consume()

		select {
		case <-consumer.done:
			log.Info("closing consumer receivedMessage channel updates")
			return
		default:
			if e != nil {
				log.Warn("consumer failed with error: ", e)
			} else {
				receivedMessages <- m
			}
		}
	}
}

func (consumer *AtLeastOnceConsumer) committer(broadcastedMessages chan kafka.Message) {
	messageBuffer := make([]kafka.Message, 0, 10)
	ticker := time.NewTicker(time.Millisecond * 500)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			consumer.commitMessages(messageBuffer...)
			for _, m := range messageBuffer {
				log.Debug("committed message with key: ", string(m.Key))
			}
			messageBuffer = messageBuffer[:0]
		case m, ok := <-broadcastedMessages:
			if !ok {
				// save progress before exiting
				consumer.commitMessages(messageBuffer...)
				log.Info("committer process closed")
				return
			}
			messageBuffer = append(messageBuffer, m)
		}
	}
}
