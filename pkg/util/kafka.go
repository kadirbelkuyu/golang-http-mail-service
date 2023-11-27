package util

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaProducer{
		Writer: w,
	}
}

func (kp *KafkaProducer) SendMessage(ctx context.Context, key, message string) error {
	return kp.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(message),
		},
	)
}
