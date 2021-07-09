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
	ctx    context.Context
}

func NewAtLeastOnceConsumer(ctx context.Context, consumerGroup string) *AtLeastOnceConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        constants.Brokers,
		GroupID:        consumerGroup,
		Topic:          constants.MessageTopic,
		IsolationLevel: kafka.ReadCommitted,
		CommitInterval: time.Second,
	})
	log.Trace("created a new consumer")
	return &AtLeastOnceConsumer{
		Reader: reader,
		ctx:    ctx,
	}
}

func (consumer *AtLeastOnceConsumer) consume() (kafka.Message, error) {
	return consumer.Reader.FetchMessage(consumer.ctx)
}

func (consumer *AtLeastOnceConsumer) commitMessages(msgs ...kafka.Message) error {
	return consumer.Reader.CommitMessages(consumer.ctx, msgs...)
}

func (consumer *AtLeastOnceConsumer) Run(receivedMessages, broadcastedMessages chan kafka.Message) {
	go consumer.committer(broadcastedMessages)
	defer func() {
		log.Info("closing consumer receivedMessage channel updates")
		close(receivedMessages)
		consumer.Reader.Close()
	}()
	log.Trace("running consumer")

	for {
		m, e := consumer.consume()

		select {
		case <-consumer.ctx.Done():
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
		// save progress before exiting
		consumer.commitMessages(messageBuffer...)
		log.Info("committer process closed")
		ticker.Stop()
	}()
	log.Trace("running committer")

	for {
		select {
		case <-ticker.C:
			consumer.commitMessages(messageBuffer...)
			for _, m := range messageBuffer {
				log.Debugf("committed message with key: '%s'", string(m.Key))
			}
			messageBuffer = messageBuffer[:0]
		case m, ok := <-broadcastedMessages:
			if !ok {
				return
			}
			messageBuffer = append(messageBuffer, m)
		}
	}
}
