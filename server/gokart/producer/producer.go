package producer

import (
	"context"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/message"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/segmentio/kafka-go"
)

type AtLeastOnceProducer struct {
	Writer *kafka.Writer
}

func NewAtLeaseOnceProducer(acks, maxAttempts int) *AtLeastOnceProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(constants.Brokers...),
		Topic:        constants.MessageTopic,
		Balancer:     kafka.CRC32Balancer{},
		MaxAttempts:  maxAttempts,
		RequiredAcks: kafka.RequiredAcks(acks),
		Async:        false,
	}
	log.Trace("created a new producer")
	return &AtLeastOnceProducer{
		Writer: w,
	}
}

func (producer *AtLeastOnceProducer) Produce(message *message.Chat) error {
	encoding, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	return producer.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(message.Topic),
			Value: encoding,
		})
}

func (producer *AtLeastOnceProducer) Close() {
	producer.Writer.Close()
	log.Trace("closed producer")
}
